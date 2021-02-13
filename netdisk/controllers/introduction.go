package controllers

import (

	beego "github.com/beego/beego/v2/server/web"
)

type IntroController struct {
	beego.Controller
}

func (c *IntroController) ShowIntro(){
	c.TplName = "introduction.html"
}
