package controllers

import (
	"bufio"
	"io"
	"log"
	"fmt"
	"netdisk/models"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"strconv"
)

func (c *MainController) ShowUpload() {
	if !c.Islogin() {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "index.html"
}

func (c *MainController) Upload() {
	file, head, err := c.GetFile("file")
	if err != nil {
		c.Ctx.WriteString("获取文件失败")
		return
	}
	defer file.Close()

	filename := head.Filename

	length := strings.Count(filename, "")

	if !CheckType(filename[length-5 : length-1]) {
		c.Ctx.WriteString("上传失败, 仅支持上传txt类型的文件")
		return
	} else if !CheckFile(filename, c.UserName()) {
		c.Ctx.WriteString("文件已存在, 请删除后重试")
		return
	}
	// ========================将文件内容存储至SGX中======================

	// 读取文件内容
    reader := bufio.NewReader(file)
	Text := ""
    for {
        str, err := reader.ReadString('\n') //读到一个换行就结束
		Text = Text + str
        if err == io.EOF {                  //io.EOF 表示文件的末尾
            break
        }
    }
	fmt.Println(Text)

	// 读取当前用户ID
	o := orm.NewOrm()
	user := models.User{Name: c.UserName()}
	err = o.Read(&user, "Name")
	if err != nil {
		c.Redirect("/login", 302)
		return
	}
	userId := strconv.Itoa(user.Id)
	fmt.Println(userId)

	// 获取文章标题，并转换成userId + ' ' + Title的形式
	filenameId := filename[0 : length-5]
	filenameId = userId + " " + filenameId
	fmt.Println(filenameId)

	// 组织RawInput
	file_update := RawInput{
		Id: filenameId,
		User: userId,
		Text: Text,
	}

	build_index_and_commit(aes_encrypt(json_to_string(file_update)))
	// ========================end of 将文件内容存储至SGX中======================

	err = c.SaveToFile("file", "fileStorage/"+c.UserName()+"/"+filename)
	log.Println(err)
	if err != nil || !InsertFile(filename, c.UserName()) {
		c.Ctx.WriteString("上传失败")
	} else {
		c.Ctx.Redirect(302, "/")
	}
}

func CheckType(filename string) bool {
	return filename == ".txt"
}

func CheckFile(filename string, userName string) bool {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.QueryTable("file").Filter("FileName", filename).Filter("UserName", userName).Values(&maps)

	return err == nil && num == 0
}

func InsertFile(filename string, userName string) bool {
	o := orm.NewOrm()
	file := models.File{}
	file.FileName = filename
	file.UserName = userName
	_, err := o.Insert(&file)

	return err == nil
}