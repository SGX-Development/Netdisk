package main

import (
	_ "NetdiskWeb/models"
	_ "NetdiskWeb/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run("0.0.0.0:10100")
}
