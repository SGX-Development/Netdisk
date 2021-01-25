package routers

import (
	"Netdisk/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/login", &controllers.MainController{},"get:Showlogin;post:Handlelogin")
    beego.Router("/register", &controllers.MainController{}, "get:ShowRegister;post:HandleRegister")
}
