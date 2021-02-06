package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"github.com/beego/beego/v2/client/orm"
	"Netdisk/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) ShowIndex(){
	o := orm.NewOrm()
	file := models.File{}
	file.Id = 1
	err := o.Read(&file, "Id")
	if err != nil {
		 c.Data["message"] = "暂无文件"
	} else {
		c.Data["filename"] = file.FileName
	}
	c.TplName = "index.html"
}

func (c *MainController) Logout() {
	status := c.GetSession("status").(UserStatus)
	if  status.islogin {
		status.islogin = false
		c.SetSession("status", status)
		c.Redirect("/login", 302)
	} else {
		c.TplName = "index.html"
		c.Data["message"] = "未登陆！"
	}
	log.Println("succeed")
	return
}

