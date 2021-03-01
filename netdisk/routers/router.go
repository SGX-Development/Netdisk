package routers

import (
	"netdisk/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:ShowIndex")

	beego.Router("/login", &controllers.MainController{}, "get:Showlogin;post:Handlelogin")
	beego.Router("/register", &controllers.MainController{}, "get:ShowRegister")
	beego.Router("/logout", &controllers.MainController{}, "get:Logout")

	beego.Router("/upload", &controllers.MainController{}, "get:ShowUpload;post:Upload")

	beego.Router("/download", &controllers.MainController{}, "get:Download")

	beego.Router("/session_key", &controllers.SessionKeyController{}, "get:SkGet;post:SkPost")

	beego.Router("/reg", &controllers.RegController{}, "get:RegGet")

	beego.Router("/delete", &controllers.DeleteController{}, "get:DeleteGet")

	beego.Router("/show", &controllers.ShowController{}, "get:ShowGet")

	beego.Router("/query", &controllers.QueryController{}, "get:QueryGet")
}
