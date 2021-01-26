package controllers

import (
	"Netdisk/models"
	"crypto/md5"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
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

	var hashtable = make(map[string]string)
	IsValid(userName, passWd, passWd, "1111@qq.com", hashtable)
	for key, value := range hashtable {
		c.Data[key] = value
	}

	if userName == "" || passWd == "" {
		c.Data["message"] = "用户名和密码不能为空"
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

	if user.Isdelete {
		c.Data["message"] = "用户不存在！"
		c.TplName = "login.html"
		return
	}

	// 检查密码是否正确
	if user.Passwd != PassWord(passWd){
		c.Data["message"] = "用户名或密码错误"
		c.TplName = "login.html"
		return
	}

	c.SetSession("islogin", true)
	c.SetSession("userName", user.Name)
	c.SetSession("userid", user.Id)

	//successfully login
	c.Ctx.Redirect(302, "http://58.196.135.54:10002")
}

// 处理注册
func (c *UserController) ShowRegister() {
	c.TplName = "register.html"
}

func (c *UserController) HandleRegister() {
	if c.GetSession("islogin")!=nil{
		c.Data["message"] = "登陆状态下不允许注册！"
		c.TplName = "register.html"
		return
	}

	userName := c.GetString("userName")
	passWd := c.GetString("passWd")
	passWd_2 := c.GetString("passWd_2")
	email := c.GetString("email")

	var hashtable = make(map[string]string)
	IsValid(userName, passWd, passWd_2, email, hashtable)
	for key, value := range hashtable {
		c.Data[key] = value
	}

	if passWd!=passWd_2 {
		c.Data["message"] = "两次密码不一致"
		c.TplName = "register.html"
		return
	}

	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Email = email
	user.Passwd = PassWord(passWd)
	user.Isactive = true
	user.Isdelete = false

	_,err := o.Insert(&user)
	if err != nil{
		c.Data["message"] = "Error"
		c.TplName = "register.html"
		return
	}

	//successfully register
	c.Ctx.Redirect(302, "http://58.196.135.54:10002/login")
}

func (c *UserController) DelAcc() {
	curSession := c.GetSession("userName")
	userName,ok := curSession.(string)
	if ok {
		user := models.User{}
		user.Name = userName
		user.Isdelete = true
		c.Redirect("/", 302)
	}
	return
}

func IsValid(userName string, passWd string, passWd2 string, email string, hashtable map[string]string) {
	valid := validation.Validation{}
	valid.Required(userName, "userName")    //userName can't be blank
	valid.Required(passWd, "passWd")        //passWd can't be blank
	valid.Required(passWd2, "passWd2")
	valid.Required(email, "email")          //email can't be blank
	valid.MaxSize(userName, 15, "userName") //userName MaxSize is 15
	valid.MinSize(userName, 3, "userName")  //userName MinSize is 3
	valid.MaxSize(passWd, 15, "passWd")     //passWd Maxsize is 15
	valid.MinSize(passWd, 6, "passWd")      //passWd Minsize is 6

	//var hashtable = make(map[string]string)

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			hashtable["Verify"+err.Key] = err.Message
		}
	}
}

func PassWord(passWd string) (pwmd5 string){
	var pwByte []byte = []byte(passWd)
	pw := md5.New()
	pw.Write(pwByte)
	cipherStr := pw.Sum(nil)
	return fmt.Sprintf("%x", cipherStr)
}