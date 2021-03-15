package controllers

import (
	// "fmt"
	// "netdisk/models"

	// "github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type IntroController struct {
	beego.Controller
}

func (c *IntroController) IntroGet() {
	if !c.Islogin() {
		c.Redirect("/login", 302)
	}

	c.Data["userName"] = c.UserName()

	c.TplName = "introduce.html"
}