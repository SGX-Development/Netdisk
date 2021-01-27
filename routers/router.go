package routers

import (
	"Netdisk/controllers"
    beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{}, "get:ShowIndex;post:Logout")
    beego.Router("/introduction", &controllers.IntroController{}, "get:ShowIntro")  //INTRODUCTION界面

    beego.Router("/login", &controllers.UserController{},"get:Showlogin;post:Handlelogin")  //登录界面
    beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")  //注册界面

    //beego.Router("/logout", &controllers.MainController{}, "get:Logout")

    beego.Router("/product", &controllers.ProductController{}, "get:Show")  //PRODUCT界面

    beego.Router("/personalcenter", &controllers.PersonalCenterController{}, "get:Show")  //PERSONAL CENTER界面

    beego.Router("/contactus", &controllers.ContactusController{}, "get:Show")  //CONTACT US界面

    beego.Router("/upload", &controllers.UploadController{}, "get:ShowUpload;post:Upload")

    beego.Router("/download", &controllers.DownloadController{}, "get:Download")
}
