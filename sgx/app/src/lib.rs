extern crate lazy_static;

extern crate sgx_types;
extern crate sgx_urts;
use sgx_types::*;
use sgx_urts::SgxEnclave;

use lazy_static::lazy_static;
use std::ptr;



extern crate crypto;
extern crate base64;
extern crate serde;


use serde::{Deserialize, Serialize};
use crypto::buffer::{BufferResult, ReadBuffer, WriteBuffer};
use crypto::{aes, blockmodes, buffer, symmetriccipher};

#[derive(Serialize, Deserialize, Debug)]
struct G {
    A: String,
}

// #[derive(Serialize, Deserialize, Debug)]
// struct Package {
//     user: String,
//     data: String,
// }


static ENCLAVE_FILE: &'static str = "enclave.signed.so";
lazy_static! {
    static ref SGX_ENCLAVE: SgxResult<SgxEnclave> = init_enclave();
}

extern "C" {
    fn build_index(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        line: *const u8,
        len: usize,
        // user: i32,
    ) -> sgx_status_t;
    fn delete_index(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        line: *const u8,
        len: usize,
    ) -> sgx_status_t;
    fn commit(eid: sgx_enclave_id_t, retval: *mut sgx_status_t) -> sgx_status_t;
    fn do_query(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        line: *const u8,
        len: usize,
        encrypted_result_string: *const u8,
        result_max_len: usize,
    ) -> sgx_status_t;

    fn get_origin_by_id(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        line: *const u8,
        len: usize,
        encrypted_result_string: *const u8,
        result_max_len: usize,
    ) -> sgx_status_t;

    fn server_hello(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        ref_tmp_pk_n: *const u8,
        len_tmp_pk_n: &mut usize,
        ref_tmp_pk_e: *const u8,
        len_tmp_pk_e: &mut usize,
        ref_tmp_certificate: *const u8,
        len_tmp_certificate: &mut usize,
        string_limit: usize,
    ) -> sgx_status_t;

    fn user_register(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        enc_user_pswd: *const u8,
        enc_user_pswd_len: usize,
        user: *const u8,
        user_len: &mut usize,
        enc_pswd: *const u8,
        enc_pswd_len: &mut usize,
        string_limit: usize,
    ) -> sgx_status_t;

    fn get_session_key(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        user: *const u8,
        user_len: usize,
        enc_sessionkey: *const u8,
        enc_sessionkey_len: usize,
    ) -> sgx_status_t;

    fn enclave_test(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
    ) -> sgx_status_t;
}

#[no_mangle]
pub extern "C" fn init_enclave() -> SgxResult<SgxEnclave> {
    println!("init_enclave");

    let mut launch_token: sgx_launch_token_t = [0; 1024];
    let mut launch_token_updated: i32 = 0;
    // call sgx_create_enclave to initialize an enclave instance
    // Debug Support: set 2nd parameter to 1
    let debug = 1;
    let mut misc_attr = sgx_misc_attribute_t {
        secs_attr: sgx_attributes_t { flags: 0, xfrm: 0 },
        misc_select: 0,
    };

    
    let x = SgxEnclave::create(
        ENCLAVE_FILE,
        debug,
        &mut launch_token,
        &mut launch_token_updated,
        &mut misc_attr,
    );
    match &x {
        Ok(r) => {
            println!("[+] Init Enclave Successful {}!", r.geteid());
      
        }
        Err(y) => {
            eprintln!("[-] Init Enclave Failed {}!", y.as_str());
        }
    };
    println!("init_enclave_finished");
    x
}

#[no_mangle]
pub extern "C" fn rust_do_query(
    some_string: *const u8,
    some_len: usize,
    result_string_limit: usize,
    // result_string: *mut u8,
    encrypted_result_string: *mut u8,
    // result_string_size: *mut usize,
    encrypted_result_string_size: *mut usize,
) -> Result<(), std::io::Error> {
    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let line = String::from_utf8(v.to_vec()).unwrap();

    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_do_query !");
            r
        }
        Err(x) => {
            eprintln!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };
    let enclave_id = enclave.geteid();

    let mut retval = sgx_status_t::SGX_SUCCESS;

    // let mut result_vec: Vec<u8> = vec![0; result_string_limit];
    // let result_slice = &mut result_vec[..];

    let mut encrypted_result_vec: Vec<u8> = vec![0; result_string_limit];
    let encrypted_result_slice = &mut encrypted_result_vec[..];

    let result = unsafe {
        do_query(
            enclave_id,
            &mut retval,
            line.as_ptr() as *const u8,
            line.len(),
            // result_slice.as_mut_ptr(),
            encrypted_result_slice.as_mut_ptr(),
            result_string_limit,
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed {}!", result.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    let mut encrypted_result_vec: Vec<u8> = encrypted_result_slice.to_vec();
    encrypted_result_vec.retain(|x| *x != 0x00u8);
    if encrypted_result_vec.len() == 0 {
        println!("emptyString");
    } else {
        let raw_result_str = String::from_utf8(encrypted_result_vec).unwrap();
        let l = raw_result_str.len();
        if l > result_string_limit {
            panic!("{} > {}", l, result_string_limit);
        }
        unsafe {
            *encrypted_result_string_size = l;
            ptr::copy_nonoverlapping(
                raw_result_str.as_ptr(),
                encrypted_result_string,
                raw_result_str.len(),
            );
        }
    }

    Ok(())
}

#[no_mangle]
pub extern "C" fn rust_build_index(
    some_string: *const u8,
    some_len: usize,
    success: *mut usize,
) -> Result<(), std::io::Error> {
    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let line = String::from_utf8(v.to_vec()).unwrap();

    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_build_index");
            r
        }
        Err(x) => {
            eprintln!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };

    // let raw_package: Package = serde_json::from_str(&line).unwrap();
    // // let package_user = raw_package.user;
    // let package_data = raw_package.data;
    // let package_user = raw_package.user.parse::<i32>().unwrap()
        
    let enclave_id = enclave.geteid();

    let mut retval = sgx_status_t::SGX_SUCCESS;

    let result = unsafe {
        build_index(
            enclave_id,
            &mut retval,
            line.as_ptr() as *const u8,
            line.len(),
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed {}!", result.as_str());
            unsafe {
                *success = 0;
            }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "build ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            unsafe {
                *success = 0;
            }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    unsafe {
        *success = 1;
    }

    Ok(())
}


#[no_mangle]
pub extern "C" fn rust_delete_index(
    some_string: *const u8,
    some_len: usize,
    success: *mut usize,
) -> Result<(), std::io::Error> {
    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let line = String::from_utf8(v.to_vec()).unwrap();
    
    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_delete_index");
            r
        }
        Err(x) => {
            eprintln!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };
    let enclave_id = enclave.geteid();

    let mut retval = sgx_status_t::SGX_SUCCESS;

    let result = unsafe {
        delete_index(
            enclave_id,
            &mut retval,
            line.as_ptr() as *const u8,
            line.len(),
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed {}!", result.as_str());
            unsafe {
                *success = 0;
            }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            unsafe {
                *success = 0;
            }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    unsafe {
        *success = 1;
    }

    Ok(())
}


#[no_mangle]
pub extern "C" fn rust_commit(success: *mut usize) -> Result<(), std::io::Error> {
    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_commit");
            r
        }
        Err(x) => {
            eprintln!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };
    let enclave_id = enclave.geteid();

    let mut retval = sgx_status_t::SGX_SUCCESS;

    let result = unsafe { commit(enclave_id, &mut retval) };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed RES {}!", result.as_str());
            unsafe {
                *success = 0;
            }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed RET {}!", retval.as_str());
            unsafe {
                *success = 0;
            }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    unsafe {
        *success = 1;
    }

    Ok(())
}

#[no_mangle]
pub extern "C" fn rust_search_title(
    some_string: *const u8,
    some_len: usize,
    result_string_limit: usize,
    // result_string: *mut u8,
    encrypted_result_string: *mut u8,
    // result_string_size: *mut usize,
    encrypted_result_string_size: *mut usize,
) -> Result<(), std::io::Error> {
    let v: &[u8] = unsafe { std::slice::from_raw_parts(some_string, some_len) };
    let line = String::from_utf8(v.to_vec()).unwrap();

    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_search_title");
            r
        }
        Err(x) => {
            println!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };
    let enclave_id = enclave.geteid();

    let mut retval = sgx_status_t::SGX_SUCCESS;

    // let mut result_vec: Vec<u8> = vec![0; result_string_limit];
    // let result_slice = &mut result_vec[..];

    let mut encrypted_result_vec: Vec<u8> = vec![0; result_string_limit];
    let encrypted_result_slice = &mut encrypted_result_vec[..];

    let result = unsafe {
        get_origin_by_id(
            enclave_id,
            &mut retval,
            line.as_ptr() as *const u8,
            line.len(),
            // result_slice.as_mut_ptr(),
            encrypted_result_slice.as_mut_ptr(),
            result_string_limit,
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed {}!", result.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    let mut encrypted_result_vec: Vec<u8> = encrypted_result_slice.to_vec();
    encrypted_result_vec.retain(|x| *x != 0x00u8);
    if encrypted_result_vec.len() == 0 {
        println!("emptyString");
    } else {
        
        let raw_result_str = String::from_utf8(encrypted_result_vec).unwrap();
        let l = raw_result_str.len();
        if l > result_string_limit {
            panic!("{} > {}", l, result_string_limit);
        }
        unsafe {
            *encrypted_result_string_size= l;
            ptr::copy_nonoverlapping(
                raw_result_str.as_ptr(),
                encrypted_result_string,
                raw_result_str.len(),
            );
        }
    
    }

    Ok(())
}

//============================================================


#[no_mangle]
pub extern "C" fn rust_server_hello(
    pk_n: *mut u8,
    pk_n_len: *mut usize,
    pk_e: *mut u8,
    pk_e_len: *mut usize,
    certificate: *mut u8,
    certificate_len: *mut usize,
    string_limit: usize,
) -> Result<(), std::io::Error> {

    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_server_hello");
            r
        }
        Err(x) => {
            println!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };

    let mut tmp_pk_n: Vec<u8> = vec![0; string_limit];
    let mut tmp_pk_e: Vec<u8> = vec![0; string_limit];
    let mut tmp_certificate: Vec<u8> = vec![0; string_limit];

    let ref_tmp_pk_n = &mut tmp_pk_n[..];
    let ref_tmp_pk_e = &mut tmp_pk_e[..];
    let ref_tmp_certificate = &mut tmp_certificate[..];


    let enclave_id = enclave.geteid();
    let mut retval = sgx_status_t::SGX_SUCCESS;

    let mut len_tmp_pk_n: usize = 0;
    let mut len_tmp_pk_e: usize = 0;
    let mut len_tmp_certificate: usize = 0;
    // let ref ref_len_tmp_pk_n = len_tmp_pk_n;
    // let mut ref_len_tmp_pk_n = 0;

    let result = unsafe {
        server_hello(
            enclave_id,
            &mut retval,
            ref_tmp_pk_n.as_mut_ptr(),
            &mut len_tmp_pk_n,
            ref_tmp_pk_e.as_mut_ptr(),
            &mut len_tmp_pk_e,
            ref_tmp_certificate.as_mut_ptr(),
            &mut len_tmp_certificate,
            string_limit,
        )
    };
    // println!("a----{}",len_tmp_pk_n);
    // println!("b----{}",len_tmp_pk_e);
    // println!("c----{}",len_tmp_certificate);

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed {}!", result.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    unsafe {
        *pk_n_len = len_tmp_pk_n;
        ptr::copy_nonoverlapping(
            ref_tmp_pk_n.as_ptr(),
            pk_n,
            *pk_n_len,
        );
        *pk_e_len = len_tmp_pk_e;
        ptr::copy_nonoverlapping(
            ref_tmp_pk_e.as_ptr(),
            pk_e,
            *pk_e_len,
        );
        *certificate_len = len_tmp_certificate;
        ptr::copy_nonoverlapping(
            ref_tmp_certificate.as_ptr(),
            certificate,
            *certificate_len,
        );
    }

    Ok(())
}


#[no_mangle]
pub extern "C" fn rust_register(
    enc_user_pswd: *const u8,
    enc_user_pswd_len: usize,
    user: *mut u8,
    user_len: *mut usize,
    enc_pswd: *mut u8,
    enc_pswd_len: *mut usize,
    success: *mut usize,
    string_limit: usize,
) -> Result<(), std::io::Error> {

    let enc_vec: &[u8] = unsafe { std::slice::from_raw_parts(enc_user_pswd, enc_user_pswd_len) };
    let enc_data = String::from_utf8(enc_vec.to_vec()).unwrap();

    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust register");
            r
        }
        Err(x) => {
            println!("[-] rust register failled {}!", x.as_str());
            unsafe{ *success = 0; }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };

    let mut user_vec: Vec<u8> = vec![0; string_limit];
    let mut enc_pswd_vec: Vec<u8> = vec![0; string_limit];

    let tmp_user = &mut user_vec[..];
    let tmp_enc_pswd = &mut enc_pswd_vec[..];

    let mut tmp_user_len: usize = 0;
    let mut tmp_enc_pswd_len: usize = 0;


    let enclave_id = enclave.geteid();
    let mut retval = sgx_status_t::SGX_SUCCESS;

    let result = unsafe {
        user_register(
            enclave_id,
            &mut retval,
            enc_data.as_ptr() as *const u8,
            enc_data.len(),
            tmp_user.as_mut_ptr(),
            &mut tmp_user_len,
            tmp_enc_pswd.as_mut_ptr(),
            &mut tmp_enc_pswd_len,
            string_limit,
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed {}!", result.as_str());
            unsafe{ *success = 0; }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            unsafe{ *success = 0; }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    unsafe {
        *user_len = tmp_user_len;
        *enc_pswd_len = tmp_enc_pswd_len;
        ptr::copy_nonoverlapping(
            tmp_user.as_ptr(),
            user,
            *user_len,
        );
        ptr::copy_nonoverlapping(
            tmp_enc_pswd.as_ptr(),
            enc_pswd,
            *enc_pswd_len,
        );
        *success = 1;
    }

    Ok(())

}



#[no_mangle]
pub extern "C" fn rust_get_session_key(
    user: *const u8,
    user_len: usize,
    enc_sessionkey: *const u8,
    enc_sessionkey_len: usize,
) -> Result<(), std::io::Error> {

    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_delete_index");
            r
        }
        Err(x) => {
            eprintln!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };
    let enclave_id = enclave.geteid();

    let mut retval = sgx_status_t::SGX_SUCCESS;

    let result = unsafe {
        get_session_key(
            enclave_id,
            &mut retval,
            user,
            user_len,
            enc_sessionkey,
            enc_sessionkey_len,
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] ECALL Enclave Failed {}!", result.as_str());
            // unsafe{ *success = 0; }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "ecall failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            // unsafe{ *success = 0; }
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }

    Ok(())
}

//============================================================


#[no_mangle]
pub extern "C" fn rust_test() -> Result<(), std::io::Error>{
    let enclave = match &*SGX_ENCLAVE {
        Ok(r) => {
            println!("[+] rust_delete_index");
            r
        }
        Err(x) => {
            eprintln!("[-] Init Enclave Failed {}!", x.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "init enclave failed",
            ));
        }
    };
    let enclave_id = enclave.geteid();
    let mut retval = sgx_status_t::SGX_SUCCESS;

    let result = unsafe {
        enclave_test(
            enclave_id,
            &mut retval,
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            eprintln!("[-] test Enclave Failed {}!", result.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "test failed",
            ));
        }
    }
    match retval {
        sgx_status_t::SGX_SUCCESS => {}
        e => {
            eprintln!("[-] ECALL Enclave Failed {}!", retval.as_str());
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ));
        }
    }
    println!("sucess!");
    Ok(())
}



//============================================================

#[no_mangle]
extern "C" fn go_encrypt(limit_length: usize, plaintext: *mut u8, plainlength: usize, ciphertext: *mut u8, cipherlength: *mut usize) 
-> Result<(), std::io::Error>{
    let v: &[u8] = unsafe { std::slice::from_raw_parts(plaintext, plainlength) };
    let line = String::from_utf8(v.to_vec()).unwrap();  

    let raw_result_str = rust_encrypt(line);
    let l = raw_result_str.len();
    if l > limit_length {
        panic!("{} > {}", l, limit_length);
    }
    unsafe {
        *cipherlength = l;
        ptr::copy_nonoverlapping(
            raw_result_str.as_ptr(),
            ciphertext,
            raw_result_str.len(),
        );
    }
    Ok(())

}

#[no_mangle]
extern "C" fn go_decrypt(limit_length: usize, ciphertext: *mut u8, cipherlength: usize, plaintext: *mut u8, plainlength: *mut usize) 
-> Result<(), std::io::Error>{
    let v: &[u8] = unsafe { std::slice::from_raw_parts(ciphertext, cipherlength) };
    let line = String::from_utf8(v.to_vec()).unwrap();  
        
    let raw_result_str = rust_decrypt(line);

    let l = raw_result_str.len();
    if l > limit_length {
        panic!("{} > {}", l, limit_length);
    }
    unsafe {
        *plainlength = l;
        ptr::copy_nonoverlapping(
            raw_result_str.as_ptr(),
            plaintext,
            raw_result_str.len(),
        );
    }
    Ok(())
}

fn rust_encrypt(message: String) -> String {
    let g: G = G {
        A: message.to_string(),
    };
    let y = serde_json::to_string(&g).unwrap();

    let mut key: [u8; 32] = [0; 32];
    let mut iv: [u8; 16] = [0; 16];

    let x: Vec<u8> = unsafe{encrypt(y.as_bytes(), &key, &iv).ok().unwrap()};

    base64::encode(&x)
}


fn encrypt(
    data: &[u8],
    key: &[u8],
    iv: &[u8],
) -> Result<Vec<u8>, symmetriccipher::SymmetricCipherError> {

    let mut encryptor =
        aes::cbc_encryptor(aes::KeySize::KeySize256, key, iv, blockmodes::PkcsPadding);

    let mut final_result = Vec::<u8>::new();
    let mut read_buffer = buffer::RefReadBuffer::new(data);
    let mut buffer = [0; 4096];
    let mut write_buffer = buffer::RefWriteBuffer::new(&mut buffer);

    loop {
        let result = encryptor.encrypt(&mut read_buffer, &mut write_buffer, true)?;
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

fn rust_decrypt(message: String) -> String {
    let ciphertext = message.as_ptr() as *const u8;
    let ciphertext_len = message.len();

    let ciphertext_slice = unsafe { std::slice::from_raw_parts(ciphertext, ciphertext_len) };

    let key: [u8; 32] = [0; 32];  //全0
    let iv: [u8; 16] = [0; 16];
    let w = base64::decode(ciphertext_slice);
  
    let z = w.unwrap();  //Vec<u8>

    let x = decrypt(&z[..], &key, &iv).unwrap();
    let y: &str = std::str::from_utf8(&x).unwrap();
    let g: G = serde_json::from_str(&y).unwrap();

    g.A
}
