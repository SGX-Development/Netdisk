package main

//#cgo LDFLAGS: -L${SRCDIR}/../sgx/app/target/release -L /opt/sgxsdk/lib64 -ltantivy -l sgx_urts -ldl -lm
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

func main() {

	// file1 := RawInput{
	// 	Id:   "Sky",
	// 	User: "1",
	// 	Text: "The sky (also sometimes called celestial dome) is everything that lies above the surface of the Earth, including the atmosphere and outer space. In the field of astronomy, the sky is also called the celestial sphere. This is an abstract sphere, concentric to the Earth, on which the Sun, Moon, planets, and stars appear to be drifting. The celestial sphere is conventionally divided into designated areas called constellations. Usually, the term sky informally refers to a perspective from the Earth's surface; however, the meaning and usage can vary. An observer on the surface of the Earth can see a small part of the sky, which resembles a dome (sometimes called the sky bowl) appearing flatter during the day than at night.[1] In some cases, such as in discussing the weather, the sky refers to only the lower, denser layers of the atmosphere. The daytime sky appears blue because air molecules scatter shorter wavelengths of sunlight more than longer ones (redder light).[2][3][4][5] The night sky appears to be a mostly dark surface or region spangled with stars. The Sun and sometimes the Moon are visible in the daytime sky unless obscured by clouds. At night, the Moon, planets, and stars are similarly visible in the sky. Some of the natural phenomena seen in the sky are clouds, rainbows, and aurorae. Lightning and precipitation are also visible in the sky. Certain birds and insects, as well as human inventions like aircraft and kites, can fly in the sky. Due to human activities, smog during the day and light pollution during the night are often seen above large cities.",
	// }

	// file2 := RawInput{
	// 	Id:   "Poetry",
	// 	User: "1",
	// 	Text: "Poetry (derived from the Greek poiesis) is a form of literature that uses aesthetic and often rhythmic[1][2][3] qualities of language—such as phonaesthetics, sound symbolism, and metre—to evoke meanings in addition to, or in place of, the prosaic ostensible meaning. Poetry has a long history – dating back to prehistoric times with hunting poetry in Africa, and to panegyric and elegiac court poetry of the empires of the Nile, Niger, and Volta River valleys.[4] Some of the earliest written poetry in Africa occurs among the Pyramid Texts written during the 25th century BCE. The earliest surviving Western Asian epic poetry, the Epic of Gilgamesh, was written in Sumerian.",
	// }

	// file3 := RawInput{
	// 	Id:   "Sea",
	// 	User: "2",
	// 	Text: "The sea is the interconnected system of all the Earth's oceanic waters, including the Atlantic, Pacific, Indian, Southern and Arctic Oceans.[1] However, the word ’sea‘ can also be used for many specific, much smaller bodies of seawater, such as the North Sea or the Red Sea. There is no sharp distinction between seas and oceans, though generally seas are smaller, and are often partly (as marginal seas) or wholly (as inland seas) bordered by land.[2] However, the Sargasso Sea has no coastline and lies within a circular current, the North Atlantic Gyre.[3](p90) Seas are generally larger than lakes and contain salt water, but the Sea of Galilee is a freshwater lake.[4][a] The United Nations Convention on the Law of the Sea states that all of the ocean is sea.",
	// }

	// file4 := RawInput{
	// 	Id:   "Summer",
	// 	User: "2",
	// 	Text: "Summer is the hottest of the four temperate seasons, falling after spring and before autumn. At or around the summer solstice (about 3 days before Midsummer Day), the earliest sunrise and latest sunset occurs, the days are longest and the nights are shortest, with day length decreasing as the season progresses after the solstice. The date of the beginning of summer varies according to climate, tradition, and culture. When it is summer in the Northern Hemisphere, it is winter in the Southern Hemisphere, and vice versa.From an astronomical view, the equinoxes and solstices would be the middle of the respective seasons,[1][2] but sometimes astronomical summer is defined as starting at the solstice, the time of maximal insolation, often identified with the 21st day of June or December. A variable seasonal lag means that the meteorological centre of the season, which is based on average temperature patterns, occurs several weeks after the time of maximal insolation.[3] The meteorological convention is to define summer as comprising the months of June, July, and August in the northern hemisphere and the months of December, January, and February in the southern hemisphere.[4][5] Under meteorological definitions, all seasons are arbitrarily set to start at the beginning of a calendar month and end at the end of a month.[4] This meteorological definition of summer also aligns with the commonly viewed notion of summer as the season with the longest (and warmest) days of the year, in which daylight predominates. The meteorological reckoning of seasons is used in Australia, Austria, Denmark, Russia and Japan. It is also used by many in the United Kingdom and in Canada. In Ireland, the summer months according to the national meteorological service, Met Éireann, are June, July and August. However, according to the Irish Calendar, summer begins on 1 May and ends on 1 August. School textbooks in Ireland follow the cultural norm of summer commencing on 1 May rather than the meteorological definition of 1 June. In midsummer, the sun can appear even at midnight in the northern hemisphere. Photo of midnight sun in Inari, Finland. Days continue to lengthen from equinox to solstice and summer days progressively shorten after the solstice, so meteorological summer encompasses the build-up to the longest day and a diminishing thereafter, with summer having many more hours of daylight than spring. Reckoning by hours of daylight alone, summer solstice marks the midpoint, not the beginning, of the seasons. Midsummer takes place over the shortest night of the year, which is the summer solstice, or on a nearby date that varies with tradition. Where a seasonal lag of half a season or more is common, reckoning based on astronomical markers is shifted half a season.[6] By this method, in North America, summer is the period from the summer solstice (usually 20 or 21 June in the Northern Hemisphere) to the autumn equinox.[7][8][9] Reckoning by cultural festivals, the summer season in the United States is traditionally regarded as beginning on Memorial Day weekend (the last Weekend in May) and ending on Labor Day (the first Monday in September), more closely in line with the meteorological definition for the parts of the country that have four-season weather. The similar Canadian tradition starts summer on Victoria Day one week prior (although summer conditions vary widely across Canada's expansive territory) and ends, as in the United States, on Labour Day. In Chinese astronomy, summer starts on or around 5 May, with the jiéqì (solar term) known as lìxià, i.e.",
	// }

	// file5 := RawInput{
	// 	Id:   "Sea",
	// 	User: "3",
	// 	Text: "Wind blowing over the surface of a body of water forms waves that are perpendicular to the direction of the wind. The friction between air and water caused by a gentle breeze on a pond causes ripples to form. A strong blow over the ocean causes larger waves as the moving air pushes against the raised ridges of water. The waves reach their maximum height when the rate at which they are travelling nearly matches the speed of the wind. In open water, when the wind blows continuously as happens in the Southern Hemisphere in the Roaring Forties, long, organised masses of water called swell roll across the ocean.[3](pp83–84)[36][37][d] If the wind dies down, the wave formation is reduced, but already-formed waves continue to travel in their original direction until they meet land. The size of the waves depends on the fetch, the distance that the wind has blown over the water and the strength and duration of that wind. When waves meet others coming from different directions, interference between the two can produce broken, irregular seas.[36] Constructive interference can cause individual (unexpected) rogue waves much higher than normal.[38] Most waves are less than 3 m (10 ft) high[38] and it is not unusual for strong storms to double or triple that height;[39] offshore construction such as wind farms and oil platforms use metocean statistics from measurements in computing the wave forces (due to for instance the hundred-year wave) they are designed against.[40] Rogue waves, however, have been documented at heights above 25 meters (82 ft).",
	// }

	// package1 := Package_to_string(Package{User: 1, Data: aes_encrypt(json_to_string(file1))})
	// package1 := Package_to_string(Package{User: 2, Data: aes_encrypt(json_to_string(file3))})

	// fmt.Println(package1)
	// package2 := Package_to_string(Package{User: 1, Data: aes_encrypt("1 sea")})

	// build_index_and_commit(package1)
	// build_index_and_commit(aes_encrypt(json_to_string(file2)))
	// build_index_and_commit(aes_encrypt(json_to_string(file3)))
	// build_index_and_commit(aes_encrypt(json_to_string(file4)))
	// build_index_and_commit(aes_encrypt(json_to_string(file5)))

	// do_query(package2)

	// package3 := Package_to_string(Package{User: 1, Data: aes_encrypt("1 Sky")})

	// delete_index_and_commit(package3)

	// do_query(aes_encrypt("1 Sky"))

	// package4 := Package_to_string(Package{User: 1, Data: aes_encrypt("1 Skaaaaaaaaaaaaaaaaaaaaaaaaaay")})
	// search_title(package4)

	// raw_pack := Package{
	// 	User: "1",
	// 	Text: "abcd",
	// }
	// pack_string := Package_to_string(raw_pack)
	// fmt.Println(pack_string)

	// pack_get := string_to_Package(pack_string)
	// fmt.Println("%+v", pack_get)

	// public_key, certificate := server_hello()
	// fmt.Println("Welcome:")
	// fmt.Println(public_key)
	// fmt.Println(certificate)
	// astring := "enc_data"
	// user, enc_pswd := user_register(astring)
	// fmt.Println(user)
	// fmt.Println(enc_pswd)
	C.rust_test()

}

//======================================================

func server_hello() (string, string) {
	pk_e := (*C.char)(C.malloc(STRING_LIMIT))
	pk_e_len := (C.ulong)(0)

	pk_n := (*C.char)(C.malloc(STRING_LIMIT))
	pk_n_len := (C.ulong)(0)

	Certificate := (*C.char)(C.malloc(STRING_LIMIT))
	Certificate_len := (C.ulong)(0)

	C.rust_server_hello(pk_n, &pk_n_len, pk_e, &pk_e_len, Certificate, &Certificate_len, STRING_LIMIT)

	public_key_n_str := C.GoStringN(pk_n, (C.int)(pk_n_len))
	fmt.Println("public_key_n_str:", public_key_n_str)
	public_key_e_str := C.GoStringN(pk_e, (C.int)(pk_e_len))
	fmt.Println("public_key_e_str:", public_key_e_str)
	Certificate_str := C.GoStringN(Certificate, (C.int)(Certificate_len))
	// fmt.Println("Certificate_str:", Certificate_str)
	pkstr := string_public_key{
		N: public_key_n_str,
		E: public_key_e_str,
	}

	// user_str := "user1"
	// get_session_key(user_str, public_key_n_str)

	publickey, err := json.Marshal(pkstr)
	if err != nil {
		panic("marshal failed")
	}
	return string(publickey), Certificate_str
}

func user_register(enc_user_pswd string) (string, string) {
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
	return false
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
