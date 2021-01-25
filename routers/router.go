package routers

import (
	"Netdisk/controllers"
    beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{}, "get:ShowIndex")
    beego.Router("/login", &controllers.UserController{},"get:Showlogin;post:Handlelogin")
    beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
}
