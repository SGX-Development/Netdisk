package controllers

import (
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"Netdisk/models"
	"github.com/beego/beego/v2/client/orm"
)

type MainController struct {
	beego.Controller
}


func (c *MainController) Showlogin() {
	c.TplName = "login.html"
}

func (c *MainController) Handlelogin() {
	userName := c.GetString("userName")
	passWd := c.GetString("passWd")
	// email := c.GetString("email")

	valid := validation.Validation{}
	valid.Required(userName, "userName")    //userName can't be blank
	valid.Required(passWd, "passWd")        //passWd can't be blank
	// valid.Required(email, "email")          //email can't be blank
	valid.MaxSize(userName, 15, "userName") //userName MaxSize is 15
	valid.MinSize(userName, 3, "userName")  //userName MinSize is 3
	valid.MaxSize(passWd, 15, "passWd")     //passWd Maxsize is 15
	valid.MinSize(passWd, 6, "passWd")      //passWd Minsize is 6

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			data := "Verify" + err.Key
			c.Data[data] = err.Message
			//log.Println(err.Key, err.Message)
		}

	}

	if userName == "" || passWd == "" {
		log.Println("输入数据不合法")
		c.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	user := models.User{}

	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil{
		log.Println("查询失败")
		c.TplName = "login.html"
		return
	}

	c.Ctx.WriteString("Welcome")
}

func (c *MainController) ShowRegister() {
	c.TplName = "register.html"
}

func (c *MainController) HandleRegister() {
	userName := c.GetString("userName")
	passWd := c.GetString("passWd")
	email := c.GetString("email")

	valid := validation.Validation{}
	valid.Required(userName, "userName")    //userName can't be blank
	valid.Required(passWd, "passWd")        //passWd can't be blank
	valid.Required(email, "email")          //email can't be blank
	valid.MaxSize(userName, 15, "userName") //userName MaxSize is 15
	valid.MinSize(userName, 3, "userName")  //userName MinSize is 3
	valid.MaxSize(passWd, 15, "passWd")     //passWd Maxsize is 15
	valid.MinSize(passWd, 6, "passWd")      //passWd Minsize is 6

	o := orm.NewOrm()

	user := models.User{}
	user.Name = userName
	user.Passwd = passWd
	user.Email = email

	_,err := o.Insert(&user)

	if err != nil{
		log.Println("插入数据库失败")
		c.Redirect("/register",302)
		return
	}

	c.TplName = "login.html"
}