package controllers

import (
	"fmt"
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

func(c *BinController) BinPost() {
	ReturnData := make(map[string]interface{})

	if empty_bin() {
		o := orm.NewOrm()
		// file := File{FileName: filename}
		_, err := o.QueryTable("file").Filter("UserName", c.UserName()).Filter("Isdelete", true).Delete()
		if err != nil {
			ReturnData["res"]='0'
			ReturnData["message"]="db delete failure"
			fmt.Println(err)
		} else {
			ReturnData["res"] = "1"
			ReturnData["message"] = "0"
		}
	} else {
		ReturnData["res"] = "0"
		ReturnData["message"] = "empty bin failure"
	}

	c.Data["json"] = ReturnData
	c.ServeJSON()
}