package controllers

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"strings"
	"Netdisk/models"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) ShowUpload() {
	curSession := c.GetSession("status")
	if curSession==nil || (curSession!=nil && !curSession.(UserStatus).islogin) {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "upload.html"
}

func (c *UploadController) Upload() {
	curSession := c.GetSession("status")
	userName := curSession.(UserStatus).userName

	file,head,err:=c.GetFile("file")
	if err!=nil {
		c.Ctx.WriteString("获取文件失败")
		return
	}
	defer file.Close()

	filename:=head.Filename

	length := strings.Count(filename, "")
	log.Println(filename)
	if filename[length-5:length-1] !=".txt" {
		c.Ctx.WriteString("上传失败, 仅支持上传txt类型的文件")
		return
	}

	//userName := c.GetSession("status").(UserStatus).userName
	//log.Println("static/"+userName+"/"+filename)

	err =c.SaveToFile("file","fileStorage/"+filename)
	log.Println(err)
	if err!=nil {
		c.Ctx.WriteString("上传失败")
	}else {
		o := orm.NewOrm()
		file := models.File{}
		file.FileName = filename
		file.UserName = userName
		_, err = o.Insert(&file)
		if err != nil {
			c.Ctx.WriteString("插入失败")
		} else {
			c.Ctx.WriteString("上传成功")
		}
	}
}


