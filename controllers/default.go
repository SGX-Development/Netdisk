package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"log"
)

type MainController struct {
	beego.Controller
}

type File struct {
	Id       int
	Name 	string
}
func (c *MainController) ShowIndex(){
	status := c.GetSession("status")
	if !Islogin(status) {
		c.Redirect("/login", 302)
	}

	o := orm.NewOrm()
	var maps []orm.Params
	num, err:= o.QueryTable("file").Filter("UserName", status.(UserStatus).userName).Values(&maps, "FileName")//.Values(&maps)

	log.Println(num)
	log.Println(maps)

	c.Data["filename"] = ""

	if err == nil {
		fmt.Printf("Result Nums: %d\n", num)
		for _, m := range maps {
			c.Data["filename"] = c.Data["filename"].(string) + m["FileName"].(string) + "\n"
			fmt.Println(m["Id"], m["FileName"])
		}
	} else {
		c.Data["message"] = "暂无文件"
	}

	c.TplName = "index.html"
}

func (c *MainController) Logout() {
	status := c.GetSession("status")
	if Islogin(status) {
		status := c.GetSession("status").(UserStatus)
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

