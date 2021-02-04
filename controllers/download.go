package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type DownloadController struct {
	beego.Controller
}


func (c *DownloadController) Download(filename string) {
	route := "fileStorage" + filename
	c.Ctx.Output.Download(route, filename)
}
