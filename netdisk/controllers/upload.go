package controllers

import (
	"fmt"
	// "log"
	// "strings"

	// "bufio"
	// "io"
	//"log"
	// "fmt"
	"netdisk/models"
	//"strings"

	"github.com/beego/beego/v2/client/orm"
	// "strconv"
)

func (c *MainController) ShowUpload() {
	if !c.Islogin() {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "main.html"
}
func (c *MainController) Upload() {
	ReturnData := make(map[string]interface{})

	package_str := c.GetString("package_str")
	filename := c.GetString("title")
	username := c.GetString("username")
	date := c.GetString("date")
	fmt.Println(package_str)
	fmt.Println(filename)
	fmt.Println(username)

	if CheckFile(filename, c.UserName()) {
		if build_index_and_commit(package_str) {
			if InsertFile(filename, username, date) {
				ReturnData["res"] = "1"
				ReturnData["message"] = "0"
			} else {
				ReturnData["res"] = "0"
				ReturnData["message"] = "文件上传失败"
			}
		} else {
			ReturnData["res"] = "0"
			ReturnData["message"] = "文件上传失败"
		}
	} else {
		ReturnData["res"] = "0"
		ReturnData["message"] = "文件已存在, 请删除后重试"
	}

	c.Data["json"] = ReturnData
	c.ServeJSON()
}

func CheckType(filename string) bool {
	return filename == ".txt"
}

func CheckFile(filename string, userName string) bool {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.QueryTable("file").Filter("FileName", filename).Filter("UserName", userName).Values(&maps)

	return err == nil && num == 0
}

func InsertFile(filename string, userName string, date string) bool {
	o := orm.NewOrm()
	file := models.File{}
	file.FileName = filename
	file.UserName = userName
	file.Date = date
	_, err := o.Insert(&file)

	return err == nil
}
