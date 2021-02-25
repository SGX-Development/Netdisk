# interface of sgx in golang


## 1 enclave大致步骤

1. server hello: send public key and certification
2. register: get encrpyted package which contains user and password. and output user and encrpyted password.
3. login: get encrypted package, which contains user, password and session key, from user. get corresponding encrypted password from database.
4. operation: input a string and get result.
5. logout: input encryted user and delete session key in enclave, return whether succeed or not.



## 2 server hello

```go
func server_hello() (string, string)
```

### 2.1 output

向外部传输(publickey, certificate)
public key 分为 N和E，在enclave内部为bigNum，传出来为十进制string。

```rust
struct RSAPublicKey {
    n: BigUint,
    e: BigUint,
}
```

## 3 register

```go
func user_register(enc_user_pswd string) (string, string)
```

### 3.1 input

用户需要向server传输，用公钥加密的(user, password)。

```rust
struct UserInfo {
    user: String,
    password: String,
}
```

### 3.2 output

enclave 返回 字符串user，被公钥加密的字符串password

## 4 login

```go
func get_session_key(enc_pswd_from_db string, enc_data string) bool
```

登录操作，

### 4.1 input

- client: 用公钥加密的(user,passworld,sessionkey)

```rust
struct SessionKeyPackage {
    user: String,
    password: String,
    key: [u8; 32],
}
```

- database: 对应的登录在数据库中的口令

### 4.2 output

bool类型，true为登陆成功，false为登陆失败

## 5 article

### 5.1 build_index_and_commit

```go
func build_index_and_commit(user_name string, enc string) bool
```

#### 5.1.1 input

* user_name：用户名(非user_id)，明文

* enc：用session key加密的，json格式

  ```go
  type RawInput struct {
  	Id   string `json:"id"`
  	User string `json:"user"` // user_id，不是user_name
  	Text string `json:"text"`
  }
  
  
  // Example
  file2 := RawInput{
    Id:   "Poetry",
    User: "1",
    Text: "Poetry is a form of literature that uses aesthetic and often rhythmic qualities of language—such as phonaesthetics, sound symbolism, and metre—to evoke meanings in addition to, or in place of, the prosaic ostensible meaning.",
  }
  ```

#### 5.1.2 output

* bool类型，上传是否成功

### 5.2 delete_index_and_commit

```go
func build_index_and_commit(user_name string, enc string) bool
```

#### 5.2.1 input

* user_name：用户名(非user_id)，明文
* enc：用session key加密的文章标题，E~p~("Poetry")

#### 5.2.2 output

* bool类型，删除是否成功

### 5.3 do_query

```go
func do_query(user_name string, enc string) string
```

#### 5.3.1 input

* user_name：用户名(非user_id)，明文
* enc：用session key加密的搜索pattern

```go
// Example
// user_id = 1, 关键词 = "hello"
// 搜索pattern = "1 hello"
```

#### 5.3.2 output

* 用session key加密的结果

```go
type Article struct {
	Id    string
	Score float32
}
```

### 5.4 search_title

```go
func search_title(user_name string, enc string) string
```

#### 5.4.1 input

* user_name：用户名(非user_id)，明文
* enc：用session key加密的文章Title

#### 5.4.2 output

* 用session key加密的结果，返回的是文章的具体内容

## 6 others

* 注意user_name与user_id的不同

* 在SGX中现有的AES加密有格式转换

  ```rust
  // Example
  // 对于"hello"，加密得到的结果是"{"A": "hello"}"
  // 对于"{"A": "hello"}"，解密得到的结果是"hello"
  ```

* SGX对传入数据的AES解密，使用原始解密

  ```go
  // 前端加密"hello"，SGX解密得到"hello"
  ```

* SGX对传出数据的AES加密，使用现有加密

  ```go
  // SGX加密"hello"，实际为对于"{"A": "hello"}"的加密，即为现有的加密形式
  ```
