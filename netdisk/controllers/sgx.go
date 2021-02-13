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
//extern void go_encrypt(size_t limit_length, char* plaintext, size_t plainlength, char* ciphertext, size_t* cipherlength);
//extern void go_decrypt(size_t limit_length, char* ciphertext, size_t cipherlength, char* plaintext, size_t* plainlength);
import "C"

// import "log"
import "encoding/json"

// import "unsafe"
import "fmt"

type RawInput struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Text string `json:"text"`
}

type Article struct {
	Id    string
	Score float32
}

type QueryRes struct {
	A []Article
}

const STRING_LIMIT = 8192

func delete_index_and_commit(input string) {
	success := (C.ulong)(0)

	C.rust_delete_index(C.CString(input), C.ulong(len(input)), &success)

	fmt.Printf("delete_index return %d\n", success)

	C.rust_commit(&success)

	fmt.Printf("commit return %d\n", success)
}

func query_all(input string) {

	userstring := (*C.char)(C.malloc(STRING_LIMIT))
	userstring_len := (C.ulong)(0)

	// C.rust_query_all(STRING_LIMIT, C.CString(input), C.ulong(len(input)), userstring, &userstring_len)

	encrypted_data := C.GoStringN(userstring, (C.int)(userstring_len))
	fmt.Println(aes_decrypt(encrypted_data))
}

func do_query(input string) {

	const result_string_limit = 4096
	a := C.CString(input)

	c_encrypted := (*C.char)(C.malloc(result_string_limit))
	d_encrypted := (C.ulong)(0)

	C.rust_do_query(a, C.ulong(len(input)), result_string_limit, c_encrypted, &d_encrypted)

	str_encrypted := C.GoStringN(c_encrypted, (C.int)(d_encrypted))
	fmt.Println(aes_decrypt(str_encrypted))
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

	C.rust_commit(&success)

	fmt.Printf("commit return %d\n", success)
}

//--------------------------------------------------

//--------------------------------------------------

func json_to_string(input RawInput) string {
	a, err := json.Marshal(input)
	fmt.Printf("a: %s\n", a)
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
