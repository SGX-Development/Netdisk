interface of sgx in golang


enclave大致步骤
1. server hello: send public key and certification
2. register: get encrpyted package which contains user and password. and output user and encrpyted password.
3. login: get encrypted package, which contains user, password and session key, from user. get corresponding encrypted password from database.
4. operation: input a string and get result.
5. logout: input encryted user and delete session key in enclave, return whether succeed or not.



## server hello
```go
func server_hello() (string, string)
```
### outout
向外部传输(publickey, certificate)
public key 分为 N和E，在enclave内部为bigNum，传出来为十进制string。
```rust
struct RSAPublicKey {
    n: BigUint,
    e: BigUint,
}
```

## register
```go
func user_register(enc_user_pswd string) (string, string)
```
### input
用户需要向server传输，用公钥加密的(user, password)。
```rust
struct UserInfo {
    user: String,
    password: String,
}
```

### output
enclave 返回 字符串user，被公钥加密的字符串password

## login
```go
func get_session_key(enc_pswd_from_db string, enc_data string) bool
```
登录操作，
### input
- client: 用公钥加密的(user,passworld,sessionkey)

```rust
struct SessionKeyPackage {
    user: String,
    password: String,
    key: [u8; 32],
}
```
- database: 对应的登录在数据库中的口令

### output
bool类型，true为登陆成功，false为登陆失败

