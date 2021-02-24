// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License..

#![crate_name = "helloworldsampleenclave"]
#![crate_type = "staticlib"]
#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]

extern crate rsa;

extern crate num_bigint;
extern crate rand;

use rand::{rngs::StdRng, SeedableRng,Rng};

use rsa::{PublicKey, RSAPrivateKey, RSAPublicKey, PaddingScheme};
use num_bigint::BigUint;
use std::collections::HashMap;

extern crate crypto;

extern crate sgx_trts;
extern crate sgx_types;
#[cfg(not(target_env = "sgx"))]
#[macro_use]
extern crate sgx_tstd as std;

extern crate base64;
extern crate lazy_static;
extern crate serde;
extern crate tantivy;

use sgx_trts::enclave;
use sgx_types::metadata::*;
use sgx_types::*;
//use sgx_trts::{is_x86_feature_detected, is_cpu_feature_supported};
use std::backtrace::{self, PrintFormat};
use std::io::{self, Write};
use std::path::{Path, PathBuf};
use std::ptr;
use std::slice;
use std::string::String;
use std::string::ToString;
use std::sync::Arc;
use std::sync::SgxRwLock as RwLock;
use std::vec::Vec;

use tantivy::collector::{Count, TopDocs};
use tantivy::merge_policy::NoMergePolicy;
use tantivy::query::QueryParser;
use tantivy::query::TermQuery;
use tantivy::schema::*;
use tantivy::tokenizer::TokenizerManager;
use tantivy::{doc, Index, IndexReader, IndexWriter, LeasedItem, ReloadPolicy, Searcher};

use serde::{Deserialize, Serialize};

use crypto::buffer::{BufferResult, ReadBuffer, WriteBuffer};
use crypto::{aes, blockmodes, buffer, symmetriccipher};
use lazy_static::lazy_static;

// #[macro_use] 
// extern crate slice_as_array;

#[derive(Serialize, Deserialize, Debug)]
struct Articles {
    A: std::vec::Vec<Article>,
}

#[derive(Serialize, Deserialize, Debug)]
struct G {
    A: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct Article {
    Id: String,
    Score: f32,
}

#[derive(Serialize, Deserialize, Debug)]
struct RawInput {
    id: String,
    user: String,
    text: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct DBInput {
    id: String,
    user: String,
    text: String,
    user_id: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct Package {
    user: i32,
    data: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct UserInfo {
    user: String,
    password: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct SessionKeyPackage {
    user: String,
    password: String,
    key: String, //string should match[u8; 32],
}

struct SessionKey {
    user: i32,
    key: [u8; 32],
    vi: [u8; 16],
}

extern crate spin;
use spin::Mutex;


lazy_static! {
    static ref keymap: Mutex<HashMap<String, [u8;32]>> = Mutex::new(HashMap::new());


    static ref schema: Schema = {
        let mut schema_builder = Schema::builder();

        schema_builder.add_text_field("id", STRING | STORED);
        schema_builder.add_text_field("user", STRING | STORED);
        schema_builder.add_text_field("text", TEXT | STORED);
        schema_builder.add_text_field("user_id", STRING | STORED);

        schema_builder.build()
    };
    static ref index: Index = {
        std::untrusted::fs::create_dir_all("idx").map_err(|e| {
            eprintln!("{}", e);
        });

        let index_path = match tantivy::directory::MmapDirectory::open(Path::new("idx")) {
            Ok(index_path) => index_path,
            Err(e) => {
                eprintln!("{}", e);
                panic!(e);
            }
        };

        let x = match Index::open_or_create(index_path, schema.clone()) {
            Ok(index1) => index1,
            Err(e) => {
                eprintln!("{}", e);
                panic!(e);
            }
        };
        x
    };
    static ref index_writer: Arc<RwLock<IndexWriter>> =
        Arc::new(RwLock::new(match index.writer(10_000_000) {
            Ok(index_writer1) => {
                index_writer1.set_merge_policy(std::boxed::Box::new(NoMergePolicy));
                index_writer1
            }
            Err(e) => {
                eprintln!("{}", e);
                panic!(e);
            }
        }));
    static ref query_parser: QueryParser = {
        let text_field = index.schema().get_field("text").expect("no all field?!");
        QueryParser::new(
            index.schema(),
            vec![text_field],
            TokenizerManager::default(),
        )
    };
    static ref reader: IndexReader = {
        match index
            .reader_builder()
            .reload_policy(ReloadPolicy::Manual)
            .try_into()
        {
            Ok(reader1) => reader1,
            Err(e) => {
                eprintln!("{}", e);
                panic!(e);
            }
        }
    };

    static ref private_key: RSAPrivateKey = {
        // let timeseed = std::time::SystemTime::now().duration_since(std::time::SystemTime::UNIX_EPOCH).unwrap();
        // let seed = timeseed.as_secs();
        // let mut rng = rand::rngs::StdRng::seed_from_u64(seed);
        let mut rng = rand::rngs::StdRng::seed_from_u64(0);
        let bits = 2048;
        RSAPrivateKey::new(&mut rng, bits).expect("failed to generate a key")
    };
    static ref public_key: RSAPublicKey = {RSAPublicKey::from(&*private_key)};
    static ref public_key_n: Vec<u8> = {(*public_key).n_to_vecu8()};
    static ref public_key_e: Vec<u8> = {(*public_key).e_to_vecu8()};

    static ref certificate: Vec<u8> = get_from_CA();

}

#[no_mangle]
pub extern "C" fn build_index(some_string: *const u8, some_len: usize) -> sgx_status_t {
    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let vraw = String::from_utf8(v.to_vec()).unwrap();  
    let package_input: Package = serde_json::from_str(&vraw).unwrap();
    let requester = package_input.user;
    let enc_data = package_input.data;

    let x = sgx_decrypt(enc_data.as_ptr() as *const u8, enc_data.len(), &requester);

    if let Err(y) = x {
        eprintln!("sgx_decrypt failed: {:?}", y);
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

    let line: String = x.unwrap();
    let raw_input: RawInput = serde_json::from_str(&line).unwrap();

    // find if user == request user in package
    let op_user = raw_input.user.clone().parse::<i32>().unwrap();
    if op_user != requester {
        eprintln!("package error");
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

    let userid = schema.get_field("user_id").unwrap();
    let input_userid = format!("{} {}", &raw_input.user.clone(), &raw_input.id.clone());
    let userid_field = Term::from_field_text(userid, &input_userid);
    let is_exist = extract_doc_given_id(&reader, &userid_field)
        .map_err(|e| {
            panic!(e);
        })
        .unwrap();

    if !is_exist.is_none() {
        println!("Build Error: article ID exists.");
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

    let db_input = DBInput{
        id: raw_input.id.clone(),
        user: raw_input.user.clone(),
        text: raw_input.text.clone(),
        user_id: input_userid,

    };
    let input_string = serde_json::to_string(&db_input).unwrap();
    // println!("line: {}", &input_string);

    let doc = match schema.parse_document(&input_string) {
        Ok(doc) => doc,
        _ => {
            return sgx_status_t::SGX_ERROR_UNEXPECTED;
        }
    };

    let index_writer_clone_1 = index_writer.clone();
    index_writer_clone_1.read().unwrap().add_document(doc);

    sgx_status_t::SGX_SUCCESS
}

#[no_mangle]
pub extern "C" fn delete_index(some_string: *const u8, some_len: usize) -> sgx_status_t {

    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let vraw = String::from_utf8(v.to_vec()).unwrap();  
    let package_input: Package = serde_json::from_str(&vraw).unwrap();
    let requester = package_input.user;
    let enc_data = package_input.data;

    let x = sgx_decrypt(enc_data.as_ptr() as *const u8, enc_data.len(), &requester);

    if let Err(y) = x {
        eprintln!("sgx_decrypt failed: {:?}", y);
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

    let line: String = x.unwrap();

    let uid = get_id_from_data(line.clone());
    if uid != requester {
        eprintln!("package error");
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

    let user_id = schema.get_field("user_id").unwrap();
    let delete_file = Term::from_field_text(user_id, &line);
    // need checking whether it exist?

    let index_writer_clone_3 = index_writer.clone();
    index_writer_clone_3.read().unwrap().delete_term(delete_file.clone());

    sgx_status_t::SGX_SUCCESS
}

#[no_mangle]
pub extern "C" fn commit() -> sgx_status_t {
    let index_writer_clone_2 = index_writer.clone();
    index_writer_clone_2.write().unwrap().commit().map_err(|e| {
        eprintln!("{}", e);
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    });
    sgx_status_t::SGX_SUCCESS
}

#[no_mangle]
pub extern "C" fn do_query(
    some_string: *const u8,
    some_len: usize,
    // result_string: *mut u8,
    encrypted_result_string: *mut u8,
    result_max_len: usize,
) -> sgx_status_t {
    get_rsa();

    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let vraw = String::from_utf8(v.to_vec()).unwrap();  
    let package_input: Package = serde_json::from_str(&vraw).unwrap();
    let requester = package_input.user;
    let enc_data = package_input.data;

    let x = sgx_decrypt(enc_data.as_ptr() as *const u8, enc_data.len(), &requester);

    if let Err(y) = x {
        eprintln!("sgx_decrypt failed: {:?}", y);
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }
    let search_pattern: String = x.unwrap();
    // println!("line: {}", line);

    // 对line进行拆分，转化为user与pattern
    let string_slice: &[u8] = unsafe { slice::from_raw_parts(search_pattern.as_ptr() as *const u8, search_pattern.len()) };
    // println!("{:?}", string_slice);
    let mut i = 0;
    for s in string_slice {
        if *s == 32 {break;}
        i += 1;
    }
    let user_id = &string_slice[0..i];
    let pattern = &string_slice[i+1..search_pattern.len()];
    let user_id = String::from_utf8(user_id.to_vec()).unwrap();
    let pattern = String::from_utf8(pattern.to_vec()).unwrap();
    // println!("{}", user);
    // println!("{}", pattern);
    let uid =user_id.clone().parse::<i32>().unwrap();

    if uid != requester {
        let error_msg = String::from("request package error");
        eprintln!("{}", error_msg);
        unsafe {
            ptr::copy_nonoverlapping(
                error_msg.as_ptr(),
                encrypted_result_string,
                error_msg.len(),
            );
        }
        return sgx_status_t::SGX_SUCCESS;
    }

    reader.reload().unwrap();
    let searcher = reader.searcher();

    let mut point = Articles { A: vec![] };

    let query = match query_parser.parse_query(&pattern) {
        Ok(query) => query,
        Err(e) => {
            eprintln!("{}", e);
            panic!(e);
        }
    };

    let top_docs = match searcher.search(&query, &TopDocs::with_limit(100)) {
        Ok(top_docs) => top_docs,
        Err(e) => {
            eprintln!("{}", e);
            panic!(e);
        }
    };

    let id = schema.get_field("id").unwrap();
    let user = schema.get_field("user").unwrap();

    for (score, doc_address) in top_docs {
        let retrieved_doc = searcher
            .doc(doc_address)
            .map_err(|e| {
                eprintln!("{}", e);
                return sgx_status_t::SGX_ERROR_UNEXPECTED;
            })
            .unwrap();

        let id = retrieved_doc.get_first(id).unwrap().text().unwrap();
        let user = retrieved_doc.get_first(user).unwrap().text().unwrap();

        if user.to_string() == user_id{
            let g = Article {
                Id: id.to_string(),
                Score: score,
            };
            // println!("{:?}", g);
            point.A.push(g);
        }
    }

    let x = serde_json::to_string(&point).unwrap();
    let encrypted_x = str2aes2base64(&x, &requester);


    if encrypted_x.len() < result_max_len {
        unsafe {
            // ptr::copy_nonoverlapping(x.as_ptr(), result_string, x.len());
            ptr::copy_nonoverlapping(
                encrypted_x.as_ptr(),
                encrypted_result_string,
                encrypted_x.len(),
            );
        }
        return sgx_status_t::SGX_SUCCESS;
    } else {
        eprintln!(
            "Result len = {} > buf size = {}",
            encrypted_x.len(),
            result_max_len
        );
        return sgx_status_t::SGX_ERROR_WASM_BUFFER_TOO_SHORT;
    }
}

#[no_mangle]
pub extern "C" fn get_origin_by_id(
    some_string: *const u8,
    some_len: usize,
    // result_string: *mut u8,
    encrypted_result_string: *mut u8,
    result_max_len: usize,
) -> sgx_status_t {

    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let vraw = String::from_utf8(v.to_vec()).unwrap();  
    let package_input: Package = serde_json::from_str(&vraw).unwrap();
    let requester = package_input.user;
    let enc_data = package_input.data;

    let x = sgx_decrypt(enc_data.as_ptr() as *const u8, enc_data.len(), &requester);

    if let Err(y) = x {
        eprintln!("sgx_decrypt failed: {:?}", y);
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

    let line: String = x.unwrap();
    println!("line: {}", &line);

    let uid = get_id_from_data(line.clone());
    if uid != requester {
        let error_msg = String::from("request package error");
        eprintln!("{}", error_msg);
        unsafe {
            ptr::copy_nonoverlapping(
                error_msg.as_ptr(),
                encrypted_result_string,
                error_msg.len(),
            );
        }
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }


    let user_id = schema.get_field("user_id").unwrap();
    let text = schema.get_field("text").unwrap();

    let frankenstein_isbn = Term::from_field_text(user_id, &line);
    let frankenstein_doc_misspelled = extract_doc_given_id(&reader, &frankenstein_isbn)
        .map_err(|e| {
            panic!(e);
        })
        .unwrap();


    if frankenstein_doc_misspelled.is_none() {
        let none_msg: &str = &String::from("None").to_string()[..];
        // eprintln!("{}", none_msg);
        let encrypted_none_msg = str2aes2base64(&none_msg, &requester);
        unsafe {
            ptr::copy_nonoverlapping(
                encrypted_none_msg.as_ptr(),
                encrypted_result_string,
                encrypted_none_msg.len(),
            );
        }
        return sgx_status_t::SGX_SUCCESS;
    }

    let y = frankenstein_doc_misspelled.unwrap();
    let x = y.get_first(text).unwrap().text().unwrap();

    let encrypted_x = str2aes2base64(&x, &requester);

    if encrypted_x.len() < result_max_len {
        unsafe {
            // ptr::copy_nonoverlapping(x.as_ptr(), result_string, x.len());
            ptr::copy_nonoverlapping(
                encrypted_x.as_ptr(),
                encrypted_result_string,
                encrypted_x.len(),
            );
        }
        return sgx_status_t::SGX_SUCCESS;
    } else {
        eprintln!(
            "Result len = {} > buf size = {}",
            encrypted_x.len(),
            result_max_len
        );
        return sgx_status_t::SGX_ERROR_WASM_BUFFER_TOO_SHORT;
    }
}

fn extract_doc_given_id(
    indexreader: &IndexReader,
    isbn_term: &Term,
) -> tantivy::Result<Option<Document>> {
    let searcher = indexreader.searcher();

    // This is the simplest query you can think of.
    // It matches all of the documents containing a specific term.
    //
    // The second argument is here to tell we don't care about decoding positions,
    // or term frequencies.
    let term_query = TermQuery::new(isbn_term.clone(), IndexRecordOption::Basic);
    let top_docs = searcher.search(&term_query, &TopDocs::with_limit(1))?;

    if let Some((_score, doc_address)) = top_docs.first() {
        let doc = searcher.doc(*doc_address)?;
        Ok(Some(doc))
    } else {
        // no doc matching this ID.
        Ok(None)
    }
}

// Decrypts a buffer with the given key and iv using
// AES-256/CBC/Pkcs encryption.
//
// This function is very similar to encrypt(), so, please reference
// comments in that function. In non-example code, if desired, it is possible to
// share much of the implementation using closures to hide the operation
// being performed. However, such code would make this example less clear.
fn decrypt(
    encrypted_data: &[u8],
    key: &[u8],
    iv: &[u8],
) -> Result<Vec<u8>, symmetriccipher::SymmetricCipherError> {
    let mut decryptor =
        aes::cbc_decryptor(aes::KeySize::KeySize256, key, iv, blockmodes::PkcsPadding);

    let mut final_result = Vec::<u8>::new();
    let mut read_buffer = buffer::RefReadBuffer::new(encrypted_data);
    let mut buffer = [0; 4096];
    let mut write_buffer = buffer::RefWriteBuffer::new(&mut buffer);

    loop {
        let result = decryptor.decrypt(&mut read_buffer, &mut write_buffer, true)?;
        final_result.extend(
            write_buffer
                .take_read_buffer()
                .take_remaining()
                .iter()
                .map(|&i| i),
        );
        match result {
            BufferResult::BufferUnderflow => break,
            BufferResult::BufferOverflow => {}
        }
    }

    Ok(final_result)
}

extern "C" fn sgx_decrypt(ciphertext: *const u8, ciphertext_len: usize, requester: &i32) -> Result<String, String> {
    let ciphertext_slice = unsafe { slice::from_raw_parts(ciphertext, ciphertext_len) };
    // println!("{:?}", ciphertext_slice);
    let key: [u8; 32] = [0; 32];
    let iv: [u8; 16] = [0; 16];
    let w = base64::decode(ciphertext_slice);
    match w {
        Err(base64::DecodeError::InvalidByte(a, b)) => {
            eprintln!("InvalidByte {} {}", a, b);
            return Err("InvalidByte".to_string());
        }
        Err(base64::DecodeError::InvalidLength) => {
            eprintln!("InvalidLength");
            return Err("InvalidLength".to_string());
        }
        Err(base64::DecodeError::InvalidLastSymbol(a, b)) => {
            eprintln!("InvalidLastSymbol {} {}", a, b);
            return Err("InvalidLastSymbol".to_string());
        }
        Ok(_) => {}
    }
    let z = w.unwrap();
    let x = decrypt(&z[..], &key, &iv).unwrap();
    let y: &str = std::str::from_utf8(&x).unwrap();
    let g: G = serde_json::from_str(&y).unwrap();
    Ok(g.A)
}

// Encrypt a buffer with the given key and iv using
// AES-256/CBC/Pkcs encryption.
fn encrypt(
    data: &[u8],
    key: &[u8],
    iv: &[u8],
) -> Result<Vec<u8>, symmetriccipher::SymmetricCipherError> {
    // Create an encryptor instance of the best performing
    // type available for the platform.
    let mut encryptor =
        aes::cbc_encryptor(aes::KeySize::KeySize256, key, iv, blockmodes::PkcsPadding);

    // Each encryption operation encrypts some data from
    // an input buffer into an output buffer. Those buffers
    // must be instances of RefReaderBuffer and RefWriteBuffer
    // (respectively) which keep track of how much data has been
    // read from or written to them.
    let mut final_result = Vec::<u8>::new();
    let mut read_buffer = buffer::RefReadBuffer::new(data);
    let mut buffer = [0; 4096];
    let mut write_buffer = buffer::RefWriteBuffer::new(&mut buffer);

    // Each encryption operation will "make progress". "Making progress"
    // is a bit loosely defined, but basically, at the end of each operation
    // either BufferUnderflow or BufferOverflow will be returned (unless
    // there was an error). If the return value is BufferUnderflow, it means
    // that the operation ended while wanting more input data. If the return
    // value is BufferOverflow, it means that the operation ended because it
    // needed more space to output data. As long as the next call to the encryption
    // operation provides the space that was requested (either more input data
    // or more output space), the operation is guaranteed to get closer to
    // completing the full operation - ie: "make progress".
    //
    // Here, we pass the data to encrypt to the enryptor along with a fixed-size
    // output buffer. The 'true' flag indicates that the end of the data that
    // is to be encrypted is included in the input buffer (which is true, since
    // the input data includes all the data to encrypt). After each call, we copy
    // any output data to our result Vec. If we get a BufferOverflow, we keep
    // going in the loop since it means that there is more work to do. We can
    // complete as soon as we get a BufferUnderflow since the encryptor is telling
    // us that it stopped processing data due to not having any more data in the
    // input buffer.
    loop {
        let result = encryptor.encrypt(&mut read_buffer, &mut write_buffer, true)?;

        // "write_buffer.take_read_buffer().take_remaining()" means:
        // from the writable buffer, create a new readable buffer which
        // contains all data that has been written, and then access all
        // of that data as a slice.
        final_result.extend(
            write_buffer
                .take_read_buffer()
                .take_remaining()
                .iter()
                .map(|&i| i),
        );

        match result {
            BufferResult::BufferUnderflow => break,
            BufferResult::BufferOverflow => {}
        }
    }

    Ok(final_result)
}

pub fn str2aes2base64(message: &str, requester: &i32) -> String {
    let g: G = G {
        A: message.to_string(),
    };
    let y = serde_json::to_string(&g).unwrap();

    let mut key: [u8; 32] = [0; 32];
    let mut iv: [u8; 16] = [0; 16];

    let x: Vec<u8> = encrypt(y.as_bytes(), &key, &iv).ok().unwrap();

    base64::encode(&x)
}


fn get_id_from_data(data: String) -> i32 {
    let string_slice: &[u8] = unsafe { slice::from_raw_parts(data.as_ptr() as *const u8, data.len()) };
    let mut i = 0;
    for s in string_slice {
        if *s == 32 {break;}
        i += 1;
    }
    let user_id = &string_slice[0..i];
    let user_id = String::from_utf8(user_id.to_vec()).unwrap();
    let uid =user_id.parse::<i32>().unwrap();
    uid
}


//---register--



fn get_from_CA() -> Vec<u8> {
    let cer = b"wo shi hao ren";
    cer.to_vec()
}

#[no_mangle]
pub extern "C" fn server_hello(
    ref_tmp_pk_n: *mut u8,
    len_tmp_pk_n: &mut usize,
    ref_tmp_pk_e: *mut u8,
    len_tmp_pk_e: &mut usize,
    ref_tmp_certificate: *mut u8,
    len_tmp_certificate: &mut usize,
    string_limit: usize,
) -> sgx_status_t {

    // match validate(&*public_key){
    //     Ok(y)=>  {
    //                 println!("[+] hello server!");
    //             }
    //     _ => {
    //         println!("[-] hello server fail!");
    //         return sgx_status_t::SGX_ERROR_UNEXPECTED;
    //     }
    // }

    println!("pk: {:?}", &*private_key);
    println!("pk: {:?}", &*public_key);
    // println!("\npkn: {:?}", &*public_key_n);
    // println!("\npke: {:?}", &*public_key_e);

    let public_key_n_str = BigUint::from_bytes_le(&*public_key_n).to_string();
    let public_key_e_str = BigUint::from_bytes_le(&*public_key_e).to_string();

    // println!("{:?}", public_key_n_str);
    // println!("{:?}", public_key_e_str);


    *len_tmp_pk_n = public_key_n_str.len() ;
    *len_tmp_pk_e = public_key_e_str.len() ;
    *len_tmp_certificate = (*certificate).len() ;

    // let m = String::from("hello");
    // let c = internals::encrypt(&public_key_n, &m);


    if ( public_key_n_str.len()< string_limit && 
        public_key_e_str.len() < string_limit &&
        (*certificate).len() < string_limit
    ) {
        unsafe {
            ptr::copy_nonoverlapping(
                public_key_n_str.as_ptr(),
                ref_tmp_pk_n,
                public_key_n_str.len(),
            );
            ptr::copy_nonoverlapping(
                public_key_e_str.as_ptr(),
                ref_tmp_pk_e,
                public_key_e_str.len(),
            );
            ptr::copy_nonoverlapping(
                (*certificate).as_ptr(),
                ref_tmp_certificate,
                (*certificate).len(),
            );
        }
    } else {
        eprintln!(
            "Public key len > buf size",
        );
        return sgx_status_t::SGX_ERROR_WASM_BUFFER_TOO_SHORT;
    }

    
    return sgx_status_t::SGX_SUCCESS;
}


#[no_mangle]
pub extern "C" fn user_register(
    enc_user_pswd: *const u8,
    enc_user_pswd_len: usize,
    user: *mut u8,
    user_len: &mut usize,
    enc_pswd: *mut u8,
    enc_pswd_len: &mut usize,
    string_limit: usize,
) -> sgx_status_t {

    println!("[+] new user!");

    let enc_vec: &[u8] = unsafe { std::slice::from_raw_parts(enc_user_pswd, enc_user_pswd_len) };
    // let enc_data = String::from_utf8(enc_vec.to_vec()).unwrap();  

    let w: &[u8] = &base64::decode(enc_vec).unwrap();
    // let enc_vec = [138, 57, 30, 230, 34, 195, 199, 159, 215, 38, 5, 169, 181, 106, 21, 203, 41, 14, 54, 76, 80, 38, 151, 11, 101, 68, 254, 221, 172, 165, 133, 231, 29, 49, 246, 73, 31, 51, 180, 221, 130, 96, 184, 40, 45, 136, 252, 246, 54, 108, 100, 248, 14, 18, 5, 158, 106, 113, 201, 26, 191, 224, 98, 159, 200, 94, 38, 176, 238, 129, 168, 211, 42, 235, 118, 119, 169, 79, 10, 51, 245, 199, 212, 190, 216, 39, 39, 206, 14, 66, 72, 171, 64, 157, 231, 84, 111, 246, 164, 0, 211, 139, 150, 204, 77, 55, 207, 186, 203, 81, 28, 6, 209, 106, 213, 196, 166, 160, 250, 88, 85, 167, 116, 113, 35, 186, 84, 170, 237, 91, 51, 199, 20, 62, 242, 176, 151, 54, 218, 79, 69, 70, 157, 83, 28, 72, 37, 155, 98, 62, 165, 106, 185, 0, 203, 245, 190, 130, 124, 207, 143, 134, 192, 8, 121, 61, 85, 71, 73, 174, 252, 219, 223, 61, 59, 188, 254, 239, 210, 57, 221, 174, 25, 247, 136, 152, 112, 118, 196, 236, 157, 219, 70, 234, 126, 168, 81, 185, 188, 63, 117, 2, 124, 36, 91, 74, 130, 217, 203, 102, 216, 167, 189, 39, 129, 150, 101, 44, 214, 138, 135, 100, 119, 140, 222, 152, 218, 226, 54, 27, 35, 161, 47, 98, 26, 28, 64, 102, 236, 245, 176, 7, 94, 185, 57, 37, 0, 255, 197, 226, 190, 227, 168, 184, 180, 200];

    let padding = PaddingScheme::new_pkcs1v15_encrypt();
    // let user_data_vec= (*private_key).decrypt(padding, enc_vec).expect("failed to decrypt");
    let user_data_vec= match (*private_key).decrypt(padding, w) {
        Ok(r) => {
            println!("[+] session key decrypt SUCCESS!");
            r
        }
        _ => {
            println!("[-] session key decrypt ERROR!");
            return sgx_status_t::SGX_ERROR_UNEXPECTED;
        }
    };

    let user_data_string = String::from_utf8(user_data_vec.to_vec()).unwrap();


    let user_data: UserInfo = serde_json::from_str(&user_data_string).unwrap();
    let tmp_user = user_data.user;
    let tmp_pswd = user_data.password;

    println!("tmp user: {}", &tmp_user);
    println!("tmp pswd: {}", &tmp_pswd);

    let mut rng = rand::rngs::StdRng::seed_from_u64(0);
    let padding = PaddingScheme::new_pkcs1v15_encrypt();
    let tmp_enc_pswd = (*public_key).encrypt(&mut rng, padding, &tmp_pswd.as_bytes()).expect("failed to encrypt");

    if ( tmp_user.len() < string_limit && tmp_user.len() < string_limit ) {
        unsafe {
            *user_len = tmp_user.len();
            *enc_pswd_len = tmp_enc_pswd.len();
            ptr::copy_nonoverlapping(
                tmp_user.as_ptr(),
                user,
                tmp_user.len(),
            );
            ptr::copy_nonoverlapping(
                tmp_enc_pswd.as_ptr(),
                enc_pswd,
                tmp_enc_pswd.len(),
            );
        }
    } else {
        eprintln!(
            "Result len > buf size",
        );
        return sgx_status_t::SGX_ERROR_WASM_BUFFER_TOO_SHORT;
    }
    
    
    return sgx_status_t::SGX_SUCCESS;
}




//-------------


#[no_mangle]
pub extern "C" fn get_session_key(
    enc_pswd_from_db: *const u8, //enc_pswd from db
    enc_pswd_from_db_len: usize,
    enc_data: *const u8, //contains user, password and session key from user.
    enc_data_len: usize,
) -> sgx_status_t {
    let enc_db_pswd_u8: &[u8] = unsafe { std::slice::from_raw_parts(enc_pswd_from_db, enc_pswd_from_db_len) };

    let enc_data_u8: &[u8] = unsafe { std::slice::from_raw_parts(enc_data, enc_data_len) };
    // println!("user: {}", &user_line);
    // println!("sk_v: {:?}", &enc_sessionkey_v);

    let db_pswd = match (*private_key).decrypt(PaddingScheme::new_pkcs1v15_encrypt(), enc_db_pswd_u8){
        Ok(r) => {
            println!("[+] password from database decrypt SUCCESS!");
            r
        }
        _ => {
            println!("[-] password from database decrypt ERROR!");
            return sgx_status_t::SGX_ERROR_UNEXPECTED;
        }
    };

    let sk_data = match (*private_key).decrypt(PaddingScheme::new_pkcs1v15_encrypt(), enc_data_u8){
        Ok(r) => {
            println!("[+] session key package decrypt SUCCESS!");
            r
        }
        _ => {
            println!("[-] session key package decrypt ERROR!");
            return sgx_status_t::SGX_ERROR_UNEXPECTED;
        }
    };

    let data_str = String::from_utf8(sk_data.to_vec()).unwrap();
    let db_pswd_str = String::from_utf8(db_pswd.to_vec()).unwrap();


    let data_struct: SessionKeyPackage  = match serde_json::from_str(&data_str){
        Ok(r) => {
            println!("[+]  package serde SUCCESS!");
            r
        }
        _ => {
            println!("[-] package serde ERROR!");
            return sgx_status_t::SGX_ERROR_UNEXPECTED;
        }
    };

    if db_pswd_str != data_struct.password {
        println!("[-] password ERROR!");
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

    let session_key = data_struct.key;
    if session_key.len() != 32 {
        println!("[-] session key length ERROR!");
        return sgx_status_t::SGX_ERROR_UNEXPECTED;
    }

      
    let mut sk:[u8;32];
    match slice_to_array_32(session_key.as_bytes().to_vec()){
        Ok(r) => {
            println!("[+] session key SUCCESS!");
            sk = r.clone();
            println!("sk: {:?}", sk);
        }
        _ =>{
            println!("[-] session key ERROR!");
            return sgx_status_t::SGX_ERROR_UNEXPECTED;
        }
    }

    (*keymap).lock().insert(data_struct.user, sk);
    println!("map: {:?}", &*keymap);


    return sgx_status_t::SGX_SUCCESS;
}

struct TryFromSliceError(());
fn slice_to_array_32<T>(slice: Vec<T>) -> Result<&'static [T; 32], TryFromSliceError> {
    if slice.len() == 32 {
        let ptr = slice.as_ptr() as *const [T; 32];
        unsafe {Ok(&*ptr)}
    } else {
        Err(TryFromSliceError(()))
    }
}


#[no_mangle]
pub extern "C" fn enclave_test() -> sgx_status_t {
    println!("[=] test SUCCESS!");

    // let test_u8:&[u8] = b"hello!";

    // println!("pk: {:?}", &*private_key);
    // println!("pk: {:?}", &*public_key);

    let getuser: UserInfo = UserInfo{
        user: String::from("take"),
        password: String::from("123456"),
    };

    // let getstr = serde_json::to_string(&getuser).unwrap();

    // println!("user string: {}", &getstr);

    // let mut rng = rand::rngs::StdRng::seed_from_u64(1);
    // let padding = PaddingScheme::new_pkcs1v15_encrypt();
    // let enc_data = (*public_key).encrypt(&mut rng, padding, getstr.as_bytes()).expect("failed to encrypt");

    // println!("user enc: {:?}", &enc_data);

    let enc_data = [138, 57, 30, 230, 34, 195, 199, 159, 215, 38, 5, 169, 181, 106, 21, 203, 41, 14, 54, 76, 80, 38, 151, 11, 101, 68, 254, 221, 172, 165, 133, 231, 29, 49, 246, 73, 31, 51, 180, 221, 130, 96, 184, 40, 45, 136, 252, 246, 54, 108, 100, 248, 14, 18, 5, 158, 106, 113, 201, 26, 191, 224, 98, 159, 200, 94, 38, 176, 238, 129, 168, 211, 42, 235, 118, 119, 169, 79, 10, 51, 245, 199, 212, 190, 216, 39, 39, 206, 14, 66, 72, 171, 64, 157, 231, 84, 111, 246, 164, 0, 211, 139, 150, 204, 77, 55, 207, 186, 203, 81, 28, 6, 209, 106, 213, 196, 166, 160, 250, 88, 85, 167, 116, 113, 35, 186, 84, 170, 237, 91, 51, 199, 20, 62, 242, 176, 151, 54, 218, 79, 69, 70, 157, 83, 28, 72, 37, 155, 98, 62, 165, 106, 185, 0, 203, 245, 190, 130, 124, 207, 143, 134, 192, 8, 121, 61, 85, 71, 73, 174, 252, 219, 223, 61, 59, 188, 254, 239, 210, 57, 221, 174, 25, 247, 136, 152, 112, 118, 196, 236, 157, 219, 70, 234, 126, 168, 81, 185, 188, 63, 117, 2, 124, 36, 91, 74, 130, 217, 203, 102, 216, 167, 189, 39, 129, 150, 101, 44, 214, 138, 135, 100, 119, 140, 222, 152, 218, 226, 54, 27, 35, 161, 47, 98, 26, 28, 64, 102, 236, 245, 176, 7, 94, 185, 57, 37, 0, 255, 197, 226, 190, 227, 168, 184, 180, 200];
    // let enc_data =  String::from_utf8(enc_vec).unwrap();

    // decryption
    let padding = PaddingScheme::new_pkcs1v15_encrypt();
    let raw_data = match (*private_key).decrypt(padding, &enc_data){
        Ok(r) => {
            println!("[+] session key decrypt SUCCESS!");
            r
        }
        _ => {
            println!("[-] session key decrypt ERROR!");
            return sgx_status_t::SGX_ERROR_UNEXPECTED;
        }
    };

    let raw_string =  String::from_utf8(raw_data.to_vec()).unwrap();

    println!("user enc: {}", raw_string);


    return sgx_status_t::SGX_SUCCESS;

}


fn get_rsa(){

    // let timeseed = std::time::SystemTime::now().duration_since(std::time::SystemTime::UNIX_EPOCH).unwrap();
    // let seed = timeseed.as_secs();
    // let mut rng = rand::rngs::StdRng::seed_from_u64(seed);
    // let bits = 2048;
    // let private_key = RSAPrivateKey::new(&mut rng, bits).expect("failed to generate a key");
    // // println!("private: {:?}", private_key);
    // let public_key = RSAPublicKey::from(&private_key);
    // println!("public: {:?}", public_key);
    // let nu8 = public_key.n_to_vecu8();
    // let eu8 = public_key.e_to_vecu8();
    // // println!("public n: {:?}", &nu8);
    // // println!("public e: {:?}", &eu8);
    // println!("public n: {}", nu8.len());
    // println!("public e: {}", eu8.len());

    // let newkey = RSAPublicKey::u8_form_pk(&nu8,&eu8);

    // // let pk_e = public_key.e().to_byte_be();
    // // let public_key_string = serde_json::to_string(&public_key).unwrap();
    // // println!("public: {}", public_key_string);


    // rng = rand::rngs::StdRng::seed_from_u64(0);

    // // Encrypt
    // let data = b"hello world";
    // // let data = data.unwrap();
    // let padding = PaddingScheme::new_pkcs1v15_encrypt();
    // let enc_data = newkey.encrypt(&mut rng, padding, &data[..]).expect("failed to encrypt");
    // // let enc_string = String::from_utf8(enc_data.clone().to_vec()).unwrap();  
    // // println!("haoye1: {:?}", enc_data);

    // // Decrypt
    // let padding = PaddingScheme::new_pkcs1v15_encrypt();
    // let dec_data = private_key.decrypt(padding, &enc_data).expect("failed to decrypt");
    // let dec_string = String::from_utf8(dec_data.clone().to_vec()).unwrap();  
    // println!("haoye2: {}", dec_string);


}
