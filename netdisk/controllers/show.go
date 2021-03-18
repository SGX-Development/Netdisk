package controllers

import (
	// "fmt"
	// "netdisk/models"

	// "errors"
	// "os"
	// "github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type ShowController struct {
	beego.Controller
}

func (c *ShowController) ShowGet() {
	ReturnData := make(map[string]interface{})

	package_str := c.GetString("package_str")
	// filename := c.GetString("title")
	// username := c.GetString("username")
	// fmt.Println(package_str)
	// fmt.Println(filename)
	// fmt.Println(username)

	ReturnData["res"] = "1"
	ReturnData["message"] = search_title(package_str)

	c.Data["json"] = ReturnData
	c.ServeJSON()

}
