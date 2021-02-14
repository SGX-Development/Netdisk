# Netdisk

## 项目结构

`netdisk`：web部分

`netdisk/controllers/sgx.go`：sgx提供的接口

`sgx`：sgx部分

`test`：原先sgx部分的main.go目录，可用于sgx部分的测试

## 项目运行

```shell
make clean
make netdisk	
cd netdisk
bee run			# 需要安装beego等
```
[beego环境配置](https://github.com/SGX-Development/GO/blob/master/beego%E7%8E%AF%E5%A2%83%E9%85%8D%E7%BD%AE.md)

### sgx 调试
在根目录下`make`编译所有文件（编译途中需要密码），`make clean`清空， `make cleandb`清空数据库。

## 传递到SGX的包
String getting from client: `user||Enc_k(data)`
user 和 data部分中的user部分必须一致，否则会报错（未全部完成）

## 已实现部分

### RawInput（上传数据）

```go
type RawInput struct {
	Id   string `json:"id"`  // Title
	User string `json:"user"`  // 用户Id
	Text string `json:"text"`  // Content
}
```

### DBInput (DB存储格式)
```rust
struct DBInput {
    id: String,
    user: String,
    text: String,
    user_id: String, // user+' '+id
}
```

### aes_encrypt

input: string (plaintext)  
output: string (encrypted text)   
rust aes 调用到Golang中的API  

### aes_decrypt

input: string (ciphertext)  
output: string (decrypted text)   
rust aes 调用到Golang中的API  
需要注意的是如果没有对应解密的明文，会报错。  

### build_index_and_commit

原来是RawInput明文传到sgx，现在改为了密文传送，并在enclave进行解密

go部分的AES256加密并未实现，当前只能手动用rust版的AES256加密得到密文（将build_index_and_commit函数第二行输出结果手动加密），再传给sgx，见build_index_and_commit函数163-173行左右（因此main.go函数中传来的数据是无效数据。。。）

### delete_index_and_commit

input: string (encrypted ID) 
output: None

删除所有与ID相同的条目。

### do_query

为了实现每个用户只能查询自己的文章，传到sgx部分的数据是“用户id + 空格 + pattern”，再对整个字符串进行加密传送

在enclave解密后进行字符串拆分，得到用户id和pattern，进行查询

关于sgx传回给后端的查询结果，只返回了加密后的结果，明文可以借助rust版AES256解密得到（这样做降低了debug的效率，传回明文的代码基本没有删，只是注释掉了，大家可以自行更改便于debug）

关于go加解密的问题同build_index_and_commit的描述

### search_title

为了实现Id（Title）的唯一性（考虑到不同用户有同名文章），所以将其格式设置为“用户Id + 空格 + 标题”（保证全局唯一），从而正确地按标题这一个属性来搜索出结果。当然可以将User和Id都传至enclave，由两个属性来搜索出结果，这样更合理更漂亮，但是我目前还不清楚Tantivy相关的代码调用

同样，也是密文传至sgx，再密文传回

关于go加解密的问题同build_index_and_commit的描述

### 文件上传的衔接部分(在没有前端加密的情况下)

后端组织成Rawinput数据结构，调用sgx函数上传即可

## TASK

* [x] 基本的上传数据，按关键字搜索，按标题搜索——已完成

* [x] go的AES-256/CBC/Pkcs算法

  参考方法：

  1. 网上找到相关的代码，只要加解密结果和rust版本一样即可（我没找到。。。）
  2. 自己实现
  3. 将enclave中加解密的代码移至sgx中，将函数接口暴露给go（go调用rust版本的加解密函数），enclave进行untrusted调用，感觉是可行的。。。但是效率我也不清楚
  4. 换一种rust和go都已经实现了的加解密算法

* [x] 实现数据的删除

  可借助Id全局唯一（见上述）这一性质来进行实现，Tantivy中有删除的example

* [ ] 用于加解密的密钥

  现在的key是32位0，iv是16位0，数据传输前双方应协商密钥，得到每个用户都不同的key（和iv），参考Diffie-Hellman算法

  关于密钥的存储等等，还没想好。。。

* [x] Makefile整合
  希望整合Makefile。在第一目录make后可以完整compile所有文件，make sgx和make netdisk可以分别compile对应目录。

* [ ] AES算法分开
  有时间可以把`app/src/lib.rs`中的AES算法分到另一个文件`app/src/aes.rs`。 

* [x] 文件上传的衔接部分(在没有前端加密的情况下)

* [ ] 文件删除的衔接部分(在没有前端加密的情况下)

* [ ] 按标题搜索的衔接部分(在没有前端加密的情况下)

* [ ] 按关键字搜索的衔接部分(在没有前端加密的情况下)

* [ ] RSA在rust下的算法没有

* [ ] search title 在搜无符合条件下有bug（未处理sgx返回值）



## 其它

[AES256_rust](https://github.com/SGX-Development/AES256_rust)

删除已构建的索引，只需删除netdisk/idx目录即可

代码中的大量注释并未删除

代码运行中的一些便于调试的打印也没有注释掉（注释掉了一些，但都没删。。。）

务必确保代码没有bug再上传至master branch，可自行创建其它branch并随意上传

大家已经实现的功能自行简单添加到“已实现部分”内容中
