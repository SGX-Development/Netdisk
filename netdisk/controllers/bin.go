package controllers

import (
	// "fmt"
	"netdisk/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type BinController struct {
	beego.Controller
}

func (c *BinController) BinGet() {
	if !c.Islogin() {
		c.Redirect("/login", 302)
	}

	var map1 = make(map[string]string)

	o := orm.NewOrm()
	var files []models.File
	o.QueryTable("file").Filter("UserName", c.UserName()).Filter("Isdelete", true).All(&files, "date", "filename")
	// o.QueryTable("file").Filter("UserName", "Emison").All(&files, "date", "filename")

	for _, file := range files {
		// fmt.Println(file.FileName)
		// fmt.Println(file.Date)
		map1[file.FileName] = file.Date
	}
	c.Data["userName"] = c.UserName()
	// fmt.Println(map1)
	c.Data["data"] = map1

	c.TplName = "main.html"
	c.TplName = "bin.html"
}
