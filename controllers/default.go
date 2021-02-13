package controllers

import (
	// "log"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

type File struct {
	Id   int
	Name string
}

func (c *MainController) ShowIndex() {
	if !c.Islogin() {
		c.Redirect("/login", 302)
	}

	c.Data["filename"], c.Data["message"] = Query(c.UserName())

	c.TplName = "index.html"
}

func Query(UserName string) (string, string) {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.QueryTable("file").Filter("UserName", UserName).Values(&maps, "FileName")

	res := ""
	message := ""
	if err == nil {
		if num == 0 {
			message = "暂无文件"
		} else {
			for _, m := range maps {
				res += m["FileName"].(string) + "\n"
			}
		}
	}

	return res, message
}
