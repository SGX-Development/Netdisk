package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"log"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) ShowUpload() {
	c.TplName = "upload.html"
}

func (c *UploadController) Upload() {
	file,head,err:=c.GetFile("file")
	if err!=nil {
		c.Ctx.WriteString("获取文件失败")
		return
	}
	defer file.Close()

	filename:=head.Filename
	userName := c.GetSession("status").(UserStatus).userName
	log.Println("static/"+userName+"/"+filename)
	err =c.SaveToFile("file","static/"+filename)
	log.Println(err)
	if err!=nil {
		c.Ctx.WriteString("上传失败")
	}else {
		c.Ctx.WriteString("上传成功")
	}
}
