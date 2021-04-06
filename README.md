# Netdisk

## 1 项目结构

`netdisk`：web部分

`netdisk/controllers/sgx.go`：sgx提供的接口

`sgx`：sgx部分

`test`：用于sgx部分的测试

## 2 项目运行

```shell
make clean
make all	
cd netdisk
bee run			# 需要安装beego等
```

[beego环境配置](https://github.com/SGX-Development/GO/blob/master/beego%E7%8E%AF%E5%A2%83%E9%85%8D%E7%BD%AE.md)