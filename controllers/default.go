package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"Netdisk/models"
	"github.com/beego/beego/v2/client/orm"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) ShowIndex(){
	c.TplName = "index.html"
}

type JsonMsg struct {
	Code int
	Msg  string
}

// 处理登录
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
		log.Println("查询失败，用户可能不存在")
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
		log.Println("密码错误")
		c.TplName = "login.html"
		return
	}
	c.SetSession("userName", user.Name)
	c.SetSession("passWd", pwmd5)

	//successfully login
	c.Ctx.Redirect(302, "http://58.196.135.54:10009")
}

/*
func (c *UserStatusController) Logout() {
	if !c.islogin {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = CreateMsg(403, "用户未登陆")
		c.ServeJSON()
		return
	}
	c.DestroySession()
	c.Data["json"] = CreateMsg(200, "用户注销成功")
	c.ServeJSON()
}
*/

func CreateMsg(code int, msg string) (m*JsonMsg){
	m = &JsonMsg{
		Code: code,
		Msg:  msg,
	}
	return
}

// 处理注册
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
		log.Println("插入数据库失败")
		c.TplName = "register.html"
		return
	}
	//successfully register
	c.Ctx.Redirect(302, "http://58.196.135.54:10009/login")
}