package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"log"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) ShowIndex(){
	c.TplName = "index.html"
}

func (c *MainController) Logout() {
	userName := c.GetSession("userName")
	if userName != nil {
		c.DelSession("userName")
		c.Redirect("/login", 302)
	} else {
		c.TplName = "index.html"
		c.Data["message"] = "未登陆！"
	}
	log.Println("succeed")
	return
}

