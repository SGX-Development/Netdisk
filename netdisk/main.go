package main

//#cgo LDFLAGS: -L${SRCDIR}/../sgx/app/target/release -L /opt/sgxsdk/lib64 -ltantivy -l sgx_urts -ldl -lm
//extern void rust_init_enclave(size_t* result);
import "C"

import (
	_ "netdisk/models"
	_ "netdisk/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if !enclave_init() {
		return
	}
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run("0.0.0.0:10011")
}

func enclave_init() bool {
	success := (C.ulong)(0)
	C.rust_init_enclave(&success)
	if success == 0 {
		return false
	} else {
		return true
	}
}
