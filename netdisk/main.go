package main

import (
	_ "netdisk/models"
	_ "netdisk/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run("0.0.0.0:10004")
}
