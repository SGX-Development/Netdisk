package controllers

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"log"
)

type MainController struct {
	beego.Controller
}

type File struct {
	Id       int
	Name 	string
}
func (c *MainController) ShowIndex(){
	status := c.GetSession("status")
	if !Islogin(status) {
		c.Redirect("/login", 302)
	}

	c.Data["filename"], c.Data["message"] = Query(UserName(status))

	c.TplName = "index.html"
}

func (c *MainController) Logout() {
	status := c.GetSession("status")
	if Islogin(status) {
		status := c.GetSession("status").(UserStatus)
		status.islogin = false
		c.SetSession("status", status)
		c.Redirect("/login", 302)
	} else {
		c.TplName = "index.html"
		c.Data["message"] = "未登陆！"
	}
	log.Println("succeed")
	return
}

func Query(UserName string) (string, string){
	o := orm.NewOrm()
	var maps []orm.Params
	num, err:= o.QueryTable("file").Filter("UserName", UserName).Values(&maps, "FileName")

	res := ""
	message := ""
	if err == nil {
		for _, m := range maps {
			res += m["FileName"].(string) + "\n"
		}
	} else if num==0 {
		message = "暂无文件"
	}

	return res, message
}
