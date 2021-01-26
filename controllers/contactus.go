package controllers

import (

	beego "github.com/beego/beego/v2/server/web"
)

type ContactusController struct {
	beego.Controller
}

func (c *ContactusController) Show(){
	c.TplName = "contactus.html"
}