package controllers

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
	"Netdisk/models"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) ShowUpload() {
	status := c.GetSession("status")
	if !Islogin(status) {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "upload.html"
}

func (c *UploadController) Upload() {
	status := c.GetSession("status")

	file,head,err:=c.GetFile("file")
	if err!=nil {
		c.Ctx.WriteString("获取文件失败")
		return
	}
	defer file.Close()

	filename:=head.Filename

	length := strings.Count(filename, "")

	if !CheckType(filename[length-5:length-1]) {
		c.Ctx.WriteString("上传失败, 仅支持上传txt类型的文件")
		return
	} else if !CheckFile(filename, UserName(status)){
		c.Ctx.WriteString("文件已存在, 请删除后重试")
		return
	}

	err =c.SaveToFile("file","fileStorage/"+filename)

	if err!=nil || !InsertFile(filename, UserName(status)) {
		c.Ctx.WriteString("上传失败")
	} else {
		c.Ctx.WriteString("上传成功")
	}
}

func CheckType(filename string) bool {
	return filename == ".txt"
}

func CheckFile(filename string, userName string) bool {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.QueryTable("file").Filter("FileName", filename).Filter("UserName", userName).Values(&maps)

	return err==nil && num==0
}

func InsertFile(filename string, userName string) bool {
	o := orm.NewOrm()
	file := models.File{}
	file.FileName = filename
	file.UserName = userName
	_, err := o.Insert(&file)

	return err == nil
}