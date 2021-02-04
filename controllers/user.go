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
//这个结构体有问题，后面修
type UserStatus struct {
	userName string
	userid int
	islogin bool
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

	errMsg, flag, id:= CheckAct(userName, passWd)
	if !flag {
		c.Data["message"] = errMsg
		c.TplName = "login.html"
		return
	}
	status := UserStatus{userName, id, true}
	c.SetSession("status", status)

	//successfully login
	c.Ctx.Redirect(302, "http://58.196.135.54:10110")
}

// 处理注册
func (c *UserController) ShowRegister() {
	c.TplName = "register.html"
}

func (c *UserController) HandleRegister() {
	curSession := c.GetSession("status")
	if curSession!=nil && curSession.(UserStatus).islogin{
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

	errMsg, flag := CheckReg(userName, passWd, passWd_2, email);
	if !flag {
		c.Data["message"] = errMsg
		c.TplName = "register.html"
		return
	}

	//successfully register
	c.Ctx.Redirect(302, "http://58.196.135.54:10110/login")
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

func CheckAct(userName string, passWd string)(errMsg string, flag bool, id int) {
	errMsg = ""
	flag = true
	id = -1
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		flag = false
		errMsg = "用户名或密码错误！"
	} else if user.Isdelete {
		flag = false
		errMsg = "用户不存在！"
	} else if user.Passwd != PassWord(passWd) {
		flag = false
		errMsg = "用户名或密码错误！"
	} else {
		id = user.Id
	}
	return errMsg, flag, id
}

func CheckReg(userName string, passWd string, passWd_2 string, email string)(errMsg string, flag bool) {
	errMsg = ""
	flag = true
	if passWd != passWd_2 {
		errMsg = "两次输入的密码不一致"
		flag = false
		return errMsg, flag
	}
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user, "Name")
	if err==nil && !user.Isdelete {
		errMsg = "该用户名已被占用！"
		flag = false
	} else {
		user.Email = email
		user.Passwd = PassWord(passWd)
		user.Isactive = true
		user.Isdelete = false
		_,err = o.Insert(&user)
		if err != nil {
			errMsg = "Error!"
			flag = false
		}
	}
	return errMsg, flag
}