package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type DownloadController struct {
	beego.Controller
}


func (c *DownloadController) Download() {
	c.Ctx.Output.Download("fileStorage/test.jpg", "test.jpg")
}
