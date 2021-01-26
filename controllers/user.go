package controllers

import (
	"Netdisk/models"
	"crypto/md5"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"log"
)

type UserController struct {
	beego.Controller
}

// 处理登录
func (c *UserController) Showlogin() {
	c.TplName = "login.html"
}

func (c *UserController) Handlelogin() {
	userName := c.GetString("userName")
	passWd := c.GetString("passWd")

	valid := validation.Validation{}
	valid.Required(userName, "userName")    //userName can't be blank
	valid.Required(passWd, "passWd")        //passWd can't be blank
	valid.MaxSize(userName, 15, "userName") //userName MaxSize is 15
	valid.MinSize(userName, 3, "userName")  //userName MinSize is 3
	valid.MaxSize(passWd, 15, "passWd")     //passWd Maxsize is 15
	valid.MinSize(passWd, 6, "passWd")      //passWd Minsize is 6

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			data := "Verify" + err.Key
			c.Data[data] = err.Message
		}
	}

	if userName == "" || passWd == "" {
		log.Println("输入数据不合法")
		c.Data["message"] = "输入数据不合法"
		c.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil{
		c.Data["message"] = "用户名或密码错误"
		c.TplName = "login.html"
		return
	}
	// 检查密码是否正确
	var pwByte []byte = []byte(passWd)
	pw := md5.New()
	pw.Write(pwByte)
	cipherStr := pw.Sum(nil)
	pwmd5 := fmt.Sprintf("%x", cipherStr)

	if user.Passwd != pwmd5{
		c.Data["message"] = "密码错误"
		c.TplName = "login.html"
		return
	}

	c.SetSession("userName", user.Name)
	c.SetSession("passWd", pwmd5)

	//successfully login
	c.Ctx.Redirect(302, "http://58.196.135.54:10011/introduction")
}

// 处理注册
func (c *UserController) ShowRegister() {
	c.TplName = "register.html"
}

func (c *UserController) HandleRegister() {
	userName := c.GetString("userName")
	passWd := c.GetString("passWd")
	passWd_2 := c.GetString("passWd_2")
	email := c.GetString("email")

	valid := validation.Validation{}
	valid.Required(userName, "userName")    //userName can't be blank
	valid.Required(passWd, "passWd")        //passWd can't be blank
	valid.Required(email, "email")          //email can't be blank
	valid.MaxSize(userName, 15, "userName") //userName MaxSize is 15
	valid.MinSize(userName, 3, "userName")  //userName MinSize is 3
	valid.MaxSize(passWd, 15, "passWd")     //passWd Maxsize is 15
	valid.MinSize(passWd, 6, "passWd")      //passWd Minsize is 6

	if passWd!=passWd_2 {
		c.Data["message"] = "两次密码不一致"
		c.TplName = "register.html"
		return
	}

	var pwByte []byte = []byte(passWd)
	pw := md5.New()
	pw.Write(pwByte)
	cipherStr := pw.Sum(nil)
	pwmd5 := fmt.Sprintf("%x", cipherStr)

	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Email = email
	user.Passwd = pwmd5
	user.Isactive = true
	user.Isdelete = false

	_,err := o.Insert(&user)

	if err != nil{
		c.Data["message"] = "Error"
		c.TplName = "register.html"
		return
	}

	//successfully register
	c.Ctx.Redirect(302, "http://58.196.135.54:10011/login")
}