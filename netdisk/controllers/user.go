package controllers

import (
	// "crypto/md5"
	"encoding/base64"
	"fmt"
	"netdisk/models"
	// "os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

type RegController struct {
	beego.Controller
}

// 处理登录
func (c *MainController) Showlogin() {
	c.TplName = "login.html"
}

func (c *MainController) Handlelogin() {
	userName := c.GetString("userName")
	passWd := c.GetString("passWd")

	var hashtable = make(map[string]string)
	IsValid(userName, passWd, "1111@qq.com", hashtable)
	for key, value := range hashtable {
		c.Data[key] = value
	}

	if userName == "" || passWd == "" {
		c.Data["message"] = "用户名和密码不能为空"
		c.TplName = "login.html"
		return
	}

	errMsg, flag, id := CheckAct(userName, passWd)
	if !flag {
		c.Data["message"] = errMsg
		c.TplName = "login.html"
		return
	}
	status := UserStatus{userName, id, true}
	c.SetSession("status", status)

	//successfully login
	c.Ctx.Redirect(302, "/introduction")
}

// 处理注册
func (c *MainController) ShowRegister() {
	c.TplName = "register.html"
}

func (c *RegController) RegGet() {

	ReturnData := make(map[string]interface{})

	// if c.Islogin() {
	// 	ReturnData["message"] = "登录状态下不允许注册！"
	// 	ReturnData["res"] = "0"
	// 	// c.Data["json"] = ReturnData
	// 	// c.ServeJSON() //响应前端
	// 	// c.StopRun()
	// }

	userName := c.GetString("userName")
	enc_uname_pwd_base64 := c.GetString("enc_uname_pwd")
	email := c.GetString("email")

	user_str, enc_pswd_str := register(enc_uname_pwd_base64)

	fmt.Println(user_str)
	fmt.Println(enc_pswd_str)

	enc_pswd_base64 := base64.StdEncoding.EncodeToString([]byte(enc_pswd_str))

	fmt.Println(enc_pswd_base64)

	// enc_pswd_base64 = enc_pswd_base64[0:255]

	var hashtable = make(map[string]string)
	IsValid(userName, enc_pswd_base64, email, hashtable)

	fmt.Println(1)

	for key, value := range hashtable {
		ReturnData[key] = value
	}

	fmt.Println(2)

	errMsg, flag := CheckReg(userName, enc_pswd_base64, email)
	if !flag {
		ReturnData["message"] = errMsg
		ReturnData["res"] = "0"
		c.Data["json"] = ReturnData
		c.ServeJSON() //响应前端
	}

	ReturnData["res"] = "1"
	ReturnData["message"] = "0"

	fmt.Println(3)
	fmt.Println(ReturnData)

	c.Data["json"] = ReturnData
	c.ServeJSON() //响应前端
}

func (c *MainController) DelAcc() {
	curSession := c.GetSession("userName")
	userName, ok := curSession.(string)
	if ok {
		user := models.User{}
		user.Name = userName
		user.Isdelete = true
		c.Redirect("/", 302)
	}
	return
}

func IsValid(userName string, passWd string, email string, hashtable map[string]string) {
	valid := validation.Validation{}
	valid.Required(userName, "userName") //userName can't be blank
	valid.Required(passWd, "passWd")     //passWd can't be blank
	// valid.Required(passWd2, "passWd2")
	valid.Required(email, "email")          //email can't be blank
	valid.MaxSize(userName, 15, "userName") //userName MaxSize is 15
	valid.MinSize(userName, 3, "userName")  //userName MinSize is 3
	// valid.MaxSize(passWd, 15, "passWd")     //passWd Maxsize is 15
	// valid.MinSize(passWd, 6, "passWd")      //passWd Minsize is 6

	//var hashtable = make(map[string]string)

	if valid.HasErrors() {
		fmt.Println("valid.HasErrors")
		for _, err := range valid.Errors {
			hashtable["Verify"+err.Key] = err.Message
		}
	}
}

// func PassWord(passWd string) (pwmd5 string) {
// 	var pwByte []byte = []byte(passWd)
// 	pw := md5.New()
// 	pw.Write(pwByte)
// 	cipherStr := pw.Sum(nil)
// 	return fmt.Sprintf("%x", cipherStr)
// }

func CheckAct(userName string, passWd string) (errMsg string, flag bool, id int) {
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
	} else if user.Passwd != passWd {
		flag = false
		errMsg = "用户名或密码错误！"
	} else {
		id = user.Id
	}

	return errMsg, flag, id
}

func CheckReg(userName string, passWd string, email string) (errMsg string, flag bool) {
	errMsg = ""
	flag = true
	// if passWd != passWd_2 {
	// 	errMsg = "两次输入的密码不一致"
	// 	flag = false
	// 	return errMsg, flag
	// }
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user, "Name")
	if err == nil && !user.Isdelete {
		errMsg = "该用户名已被占用！"
		flag = false
	} else {
		user.Email = email
		user.Passwd = passWd[0:255]
		user.Passwd_more = passWd[255:]
		fmt.Println(user.Passwd)
		fmt.Println(user.Passwd_more)
		user.Isactive = true
		user.Isdelete = false
		_, err = o.Insert(&user)
		fmt.Println(err)
		if err != nil {
			errMsg = "Error!"
			flag = false
		}
	}
	fmt.Println(errMsg)
	fmt.Println(flag)
	return errMsg, flag
}

// 处理退出
func (c *MainController) Logout() {
	if c.Islogin() {
		status := c.GetSession("status").(UserStatus)
		status.islogin = false
		c.SetSession("status", status)
	} else {
		c.Data["message"] = "未登录！"
	}

	c.Redirect("/login", 302)
}
