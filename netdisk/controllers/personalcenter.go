package controllers

import (

	beego "github.com/beego/beego/v2/server/web"
)

type PersonalCenterController struct {
	beego.Controller
}

func (c *PersonalCenterController) Show(){
	c.TplName = "personalcenter.html"
}