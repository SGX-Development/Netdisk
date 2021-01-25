package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id int
	Name string `orm:"unique"`
	Passwd string
	Email string
}

func init(){
	orm.RegisterModel(new(User))
}
