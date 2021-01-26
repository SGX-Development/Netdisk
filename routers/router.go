package routers

import (
	"Netdisk/controllers"
    beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/introduction", &controllers.MainController{}, "get:ShowIndex")  //INTRODUCTION界面

    beego.Router("/login", &controllers.UserController{},"get:Showlogin;post:Handlelogin")  //登录界面
    beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")  //注册界面

    beego.Router("/product", &controllers.ProductController{}, "get:Show")  //PRODUCT界面

    beego.Router("/personalcenter", &controllers.PersonalCenterController{}, "get:Show")  //PERSONAL CENTER界面

    beego.Router("/contactus", &controllers.ContactusController{}, "get:Show")  //CONTACT US界面



}
