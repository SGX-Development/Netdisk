package controllers

import (
	// "fmt"
	// "netdisk/models"

	// "github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type RecordController struct {
	beego.Controller
}

func (c *RecordController) RecordGet() {
	if !c.Islogin() {
		c.Redirect("/login", 302)
	}

	c.Data["userName"] = c.UserName()

	c.TplName = "record.html"
}