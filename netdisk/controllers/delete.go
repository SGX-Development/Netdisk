package controllers

import (
	// "fmt"
	// "netdisk/models"

	// "errors"
	// "os"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type DeleteController struct {
	beego.Controller
}

func (c *DeleteController) DeleteGet() {
	ReturnData := make(map[string]interface{})

	package_str := c.GetString("package_str")
	filename := c.GetString("title")
	username := c.GetString("username")
	// fmt.Println(package_str)
	// fmt.Println(filename)
	// fmt.Println(username)

	if delete_index_and_commit(package_str) {
		o := orm.NewOrm()
		// file := File{FileName: filename}
		o.QueryTable("file").Filter("UserName", username).Filter("FileName", filename).Update(orm.Params{
			"Isdelete": true,
		})
		ReturnData["res"] = "1"
		ReturnData["message"] = "0"
	} else {
		ReturnData["res"] = "0"
		ReturnData["message"] = "delete failure"
	}

	c.Data["json"] = ReturnData
	c.ServeJSON()

}
