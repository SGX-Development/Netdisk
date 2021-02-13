package controllers

import (

	beego "github.com/beego/beego/v2/server/web"
)

type ContactusController struct {
	beego.Controller
}

func (c *ContactusController) Show(){
	// do_query(aes_encrypt("3 the"))
	// search_title(aes_encrypt("1 Sky"))
	c.TplName = "contactus.html"
}