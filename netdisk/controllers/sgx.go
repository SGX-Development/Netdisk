package controllers

//#cgo LDFLAGS: -L${SRCDIR}/../../sgx/app/target/release -L /opt/sgxsdk/lib64 -ltantivy -l sgx_urts -ldl -lm
//#include <stdint.h>
//#include <math.h>
//extern unsigned long long init_enclave();
//extern void rust_do_query( char* some_string, size_t some_len,size_t result_string_limit,char * encrypted_result_string,size_t * encrypted_result_string_size);
//extern void rust_build_index( char* some_string, size_t some_len,size_t * result_string_size);
//extern void rust_delete_index( char* some_string, size_t some_len,size_t * result_string_size);
//extern void rust_search_title( char * some_string, size_t some_len,size_t result_string_limit, char * encrypted_result_string,size_t * encrypted_result_string_size);
//extern void rust_commit(size_t* result);
//extern void rust_server_hello( char* pk_n, size_t* pk_n_len, char* pk_e, size_t* pk_e_len, char* certificate, size_t* certificate_len, size_t string_limit);
//extern void rust_get_session_key(char* enc_pswd_from_db, size_t enc_pswd_from_db_len, char* enc_data, size_t enc_data_len, size_t* result_string_size);
//extern void rust_register(char* enc_user_pswd, size_t enc_user_pswd_len, char* user, size_t* user_len, char* enc_pswd, size_t* enc_pswd_len, size_t* result_string_size, size_t string_limit);
//extern void rust_user_logout( char* some_string, size_t some_len,size_t* result_string_size);
//extern void go_encrypt(size_t limit_length, char* plaintext, size_t plainlength, char* ciphertext, size_t* cipherlength);
//extern void go_decrypt(size_t limit_length, char* ciphertext, size_t cipherlength, char* plaintext, size_t* plainlength);
//extern void rust_test();
import "C"

// import "log"
import (
	"encoding/json"
	"fmt"
)

// import "unsafe"

type RawInput struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Text string `json:"text"`
}

type Package struct { // package from front
	User int32  `json:"user"`
	Data string `json:"data"`
}

type Article struct {
	Id    string
	Score float32
}

type QueryRes struct {
	A []Article
}

type string_public_key struct {
	N string `json:"n"`
	E string `json:"e"`
}

const STRING_LIMIT = 4096

//======================================================

func server_hello() (string, string, string) {
	pk_e := (*C.char)(C.malloc(STRING_LIMIT))
	pk_e_len := (C.ulong)(0)

	pk_n := (*C.char)(C.malloc(STRING_LIMIT))
	pk_n_len := (C.ulong)(0)

	Certificate := (*C.char)(C.malloc(STRING_LIMIT))
	Certificate_len := (C.ulong)(0)

	C.rust_server_hello(pk_n, &pk_n_len, pk_e, &pk_e_len, Certificate, &Certificate_len, STRING_LIMIT)

	public_key_n_str := C.GoStringN(pk_n, (C.int)(pk_n_len))
	// fmt.Println("public_key_n_str:", public_key_n_str)
	public_key_e_str := C.GoStringN(pk_e, (C.int)(pk_e_len))
	// fmt.Println("public_key_e_str:", public_key_e_str)
	Certificate_str := C.GoStringN(Certificate, (C.int)(Certificate_len))
	// fmt.Println("Certificate_str:", Certificate_str)
	// pkstr := string_public_key{
	// 	N: public_key_n_str,
	// 	E: public_key_e_str,
	// }

	// user_str := "user1"
	// get_session_key(user_str, public_key_n_str)

	// publickey, err := json.Marshal(pkstr)
	// if err != nil {
	// 	panic("marshal failed")
	// }
	return public_key_n_str, public_key_e_str, Certificate_str
}

func register(enc_user_pswd string) (string, string) {
	enc_pswd := (*C.char)(C.malloc(STRING_LIMIT))
	enc_pswd_len := (C.ulong)(0)

	user := (*C.char)(C.malloc(STRING_LIMIT))
	user_len := (C.ulong)(0)

	success := (C.ulong)(0)

	C.rust_register(C.CString(enc_user_pswd), C.ulong(len(enc_user_pswd)),
		user, &user_len, enc_pswd, &enc_pswd_len, &success, STRING_LIMIT)

	if success == 0 {
		fmt.Println("Register Failed!")
		return "", ""
	}

	user_str := C.GoStringN(user, (C.int)(user_len))
	enc_pswd_str := C.GoStringN(enc_pswd, (C.int)(enc_pswd_len))
	return user_str, enc_pswd_str
}

func get_session_key(enc_pswd_from_db string, enc_data string) bool {
	success := (C.ulong)(0)
	C.rust_get_session_key(
		C.CString(enc_pswd_from_db), C.ulong(len(enc_pswd_from_db)),
		C.CString(enc_data), C.ulong(len(enc_data)),
		&success,
	)
	if success == 0 {
		fmt.Println("session key process Failed!")
		return false
	} else {
		fmt.Println("session key process SUCCESS!")
		return true
	}
}

func user_logout(input string) bool {

	success := (C.ulong)(0)
	C.rust_user_logout(C.CString(input), C.ulong(len(input)), &success)

	// fmt.Printf("logout return %d\n", success)
	if success == 1 {
		return true
	} else {
		return false
	}

}

// ============================================

func delete_index_and_commit(input string) {

	success := (C.ulong)(0)
	fmt.Println("delete_index")
	C.rust_delete_index(C.CString(input), C.ulong(len(input)), &success)

	fmt.Printf("delete_index return %d\n", success)

	C.rust_commit(&success)

	fmt.Printf("commit return %d\n", success)
}

// func query_all(input string) {

// 	userstring := (*C.char)(C.malloc(STRING_LIMIT))
// 	userstring_len := (C.ulong)(0)

// 	// C.rust_query_all(STRING_LIMIT, C.CString(input), C.ulong(len(input)), userstring, &userstring_len)

// 	encrypted_data := C.GoStringN(userstring, (C.int)(userstring_len))
// 	fmt.Println(aes_decrypt(encrypted_data))
// }

func do_query(input string) {

	const result_string_limit = 4096
	a := C.CString(input)

	c_encrypted := (*C.char)(C.malloc(result_string_limit))
	d_encrypted := (C.ulong)(0)

	C.rust_do_query(a, C.ulong(len(input)), result_string_limit, c_encrypted, &d_encrypted)

	str_encrypted := C.GoStringN(c_encrypted, (C.int)(d_encrypted))
	fmt.Println(aes_decrypt(str_encrypted))
	fmt.Println("query done!")

}

func search_title(title string) {

	const result_string_limit = 4096
	a := C.CString(title)

	c_encrypted := (*C.char)(C.malloc(result_string_limit))
	d_encrypted := (C.ulong)(0)

	C.rust_search_title(a, C.ulong(len(title)), result_string_limit, c_encrypted, &d_encrypted)

	str_encrypted := C.GoStringN(c_encrypted, (C.int)(d_encrypted))
	fmt.Println(aes_decrypt(str_encrypted))

}

//实际上就是上传一条数据
func build_index_and_commit(input string) {

	success := (C.ulong)(0)

	C.rust_build_index(C.CString(input), C.ulong(len(input)), &success)

	fmt.Printf("build_index return %d\n", success)

	// if success == 0 {
	// 	return
	// }

	C.rust_commit(&success)

	fmt.Printf("commit return %d\n", success)
}

//--------------------------------------------------

//--------------------------------------------------

func string_to_Package(input string) Package {
	str := []byte(input)

	pack := Package{}
	err := json.Unmarshal(str, &pack)

	if err != nil {
		panic("unmarshal failed")
	}

	fmt.Println("%+v", pack)
	return pack
}

func Package_to_string(input Package) string {
	a, err := json.Marshal(input)
	// fmt.Printf("a: %s\n", a)
	if err != nil {
		panic("marshal failed")
	}
	return string(a)
}

func json_to_string(input RawInput) string {
	a, err := json.Marshal(input)
	// fmt.Printf("a: %s\n", a)
	if err != nil {
		panic("marshal failed")
	}
	return string(a)
}

func aes_encrypt(input string) string {
	cipher_t := (*C.char)(C.malloc(STRING_LIMIT))
	cipher_l := (C.ulong)(0)

	C.go_encrypt(STRING_LIMIT, C.CString(input), C.ulong(len(input)), cipher_t, &cipher_l)
	ciphertext := C.GoStringN(cipher_t, (C.int)(cipher_l))
	return ciphertext
}

func aes_decrypt(input string) string {
	plain_t := (*C.char)(C.malloc(STRING_LIMIT))
	plain_l := (C.ulong)(0)

	C.go_decrypt(STRING_LIMIT, C.CString(input), C.ulong(len(input)), plain_t, &plain_l)
	plaintext := C.GoStringN(plain_t, (C.int)(plain_l))
	return plaintext
}
