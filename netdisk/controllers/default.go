package controllers

import (
	// "log"
	"fmt"
	"netdisk/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

type File struct {
	Id   int
	Name string
}

func (c *MainController) ShowIndex() {
	if !c.Islogin() {
		c.Redirect("/login", 302)
	}

	var map1 = make(map[string]string)

	o := orm.NewOrm()
	var files []models.File
	o.QueryTable("file").Filter("UserName", c.UserName()).All(&files, "date", "filename")
	// o.QueryTable("file").Filter("UserName", "Emison").All(&files, "date", "filename")

	for _, file := range files {
		fmt.Println(file.FileName)
		fmt.Println(file.Date)
		map1[file.FileName] = file.Date
	}
	fmt.Println(map1)
	c.Data["data"] = map1

	c.TplName = "index.html"
}

func Query(UserName string) (string, string) {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.QueryTable("file").Filter("UserName", UserName).Values(&maps, "FileName")

	res := ""
	message := ""
	if err == nil {
		if num == 0 {
			message = "暂无文件"
		} else {
			for _, m := range maps {
				res += m["FileName"].(string) + "\n"
			}
		}
	}

	return res, message
}
