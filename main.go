package main

import (
	_ "Netdisk/routers"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:sgx12345@/demo1?charset=utf8")
}

func main() {
	beego.Run("0.0.0.0:10003")
}

