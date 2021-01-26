package controllers

import (

	beego "github.com/beego/beego/v2/server/web"
)

type ProductController struct {
	beego.Controller
}

func (c *ProductController) Show(){
	c.TplName = "product.html"
}