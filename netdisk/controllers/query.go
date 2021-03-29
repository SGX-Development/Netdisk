package controllers

import (
	// "fmt"
	// "netdisk/models"

	// "errors"
	// "os"
	// "github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type QueryController struct {
	beego.Controller
}

func (c *QueryController) QueryGet() {
	ReturnData := make(map[string]interface{})

	package_str := c.GetString("package_str")
	isfuzzy := c.GetString("isfuzzy")

	if(isfuzzy == "1"){
		ReturnData["message"] = do_query(package_str, 1)
	} else {
		ReturnData["message"] = do_query(package_str, 0)
	}
	// fmt.Println(package_str)

	ReturnData["res"] = "1"


	c.Data["json"] = ReturnData
	c.ServeJSON()

}
