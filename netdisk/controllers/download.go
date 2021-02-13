package controllers

import (
	"log"
)

func (c *MainController) Download() {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	//route := "fileStorage" + filename
	route := "fileStorage/" + c.UserName() + "/" + "test.txt"
	log.Println(route)
	c.Ctx.Output.Download(route)
}
/*
func (c *DownloadController) Download(filename string) {
	route := "fileStorage" + filename
	c.Ctx.Output.Download(route, filename)
}
*/