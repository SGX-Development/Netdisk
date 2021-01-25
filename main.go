package main

import (
	_ "Netdisk/routers"
	_ "Netdisk/models"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run("0.0.0.0:10008")
}
