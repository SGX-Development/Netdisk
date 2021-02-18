package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

type ContactusController struct {
	beego.Controller
}

func (c *ContactusController) Show(){
	public_key, certificate := server_hello()
	fmt.Println("Welcome:")
	fmt.Println(public_key)
	fmt.Println(certificate)
	c.TplName = "contactus.html"
}