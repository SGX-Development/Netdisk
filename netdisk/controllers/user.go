package controllers

import (
	"encoding/base64"
	"fmt"
	"netdisk/models"

	// "errors"
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
	ReturnData := make(map[string]interface{})

	userName := c.GetString("userName")
	enc_session_package_base64 := c.GetString("enc_session_package")

	enc_session_package, _ := base64.StdEncoding.DecodeString(enc_session_package_base64)

	// fmt.Println(enc_session_package)

	errMsg := ""
	id := -1
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		errMsg = "用户名不存在！"
	} else if user.Isdelete {
		errMsg = "用户已注销！"
	} else {
		id = user.Id
	}

	// fmt.Println(errMsg)
	// fmt.Println(id)

	pswd_base64 := user.Passwd + user.Passwd_more
	// fmt.Println(pswd_base64)

	pswd, _ := base64.StdEncoding.DecodeString(pswd_base64)
	// fmt.Println("hi")
	if errMsg == "" {
		if !get_session_key(string(pswd[:]), string(enc_session_package[:])) {
			errMsg = "密码错误！"
		}
	}

	status := UserStatus{userName, id, true}
	c.SetSession("status", status)

	if errMsg == "" {
		ReturnData["res"] = "1"
		ReturnData["message"] = "0"
	} else {
		ReturnData["res"] = "0"
		ReturnData["message"] = errMsg
	}

	c.Data["json"] = ReturnData
	c.ServeJSON()
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
	// 	c.Data["json"] = ReturnData
	// 	c.ServeJSON()
	// }

	userName := c.GetString("userName")
	enc_uname_pwd_base64 := c.GetString("enc_uname_pwd")
	email := c.GetString("email")

	_, enc_pswd_str := register(enc_uname_pwd_base64)
	// fmt.Println("enc_pswd_str: ", enc_pswd_str)

	enc_pswd_base64 := base64.StdEncoding.EncodeToString([]byte(enc_pswd_str))
	// fmt.Println("enc_pswd_base64: ", enc_pswd_base64)

	// var hashtable = make(map[string]string)
	// IsValid(userName, enc_pswd_base64, email, hashtable)

	// for key, value := range hashtable {
	// 	ReturnData[key] = value
	// }

	errMsg, flag := CheckReg(userName, enc_pswd_base64, email)
	if !flag {
		ReturnData["message"] = errMsg
		ReturnData["res"] = "0"
		c.Data["json"] = ReturnData
		c.ServeJSON()
	}

	ReturnData["res"] = "1"
	ReturnData["message"] = "0"

	c.Data["json"] = ReturnData
	c.ServeJSON()
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
	valid.Required(userName, "userName")    //userName can't be blank
	valid.Required(passWd, "passWd")        //passWd can't be blank
	valid.Required(email, "email")          //email can't be blank
	valid.MaxSize(userName, 15, "userName") //userName MaxSize is 15
	valid.MinSize(userName, 3, "userName")  //userName MinSize is 3
	// valid.MaxSize(passWd, 15, "passWd")     //passWd Maxsize is 15
	// valid.MinSize(passWd, 6, "passWd")      //passWd Minsize is 6

	//var hashtable = make(map[string]string)

	if valid.HasErrors() {
		// fmt.Println("valid.HasErrors")
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
		// fmt.Println(user.Passwd)
		// fmt.Println(user.Passwd_more)
		user.Isactive = true
		user.Isdelete = false
		_, err = o.Insert(&user)
		// fmt.Println(err)
		if err != nil {
			errMsg = fmt.Sprint(err)
			flag = false
		}
	}
	return errMsg, flag
}

// 处理退出
func (c *MainController) Logout() {
	ReturnData := make(map[string]interface{})

	package_str := c.GetString("package_str")
	// fmt.Println(package_str)

	// fmt.Println(user_logout(package_str))

	if c.Islogin() {
		status := c.GetSession("status").(UserStatus)
		status.islogin = false
		c.SetSession("status", status)
		if user_logout(package_str) {
			ReturnData["res"] = "1"
			ReturnData["message"] = "logout success"
		} else {
			ReturnData["res"] = "0"
			ReturnData["message"] = "logout faliure"
		}
	} else {
		ReturnData["res"] = "0"
		ReturnData["message"] = "未登录！"
	}

	// if user_logout(package_str) {
	// 	ReturnData["res"] = "1"
	// } else {
	// 	ReturnData["res"] = "0"
	// }

	c.Data["json"] = ReturnData
	c.ServeJSON()
}
