package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id int `json:"id"`
	Name string `json:"name" gorm:"type:varchar(45) not null;unique"`
	Passwd string `json:"password" gorm:"type:varchar(255)"`
	Email string `json:"email" orm:"unique"`
	Iconpath string  `json:"iconpath" gorm:"type:varchar(512);null"`
	Isactive bool  `json:"isactive" gorm:"default:true"`
	Isdelete bool  `json:"isdelete" gorm:"default:false"`
}

func init(){
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:sgx12345@/netdisk?charset=utf8")
	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, true)
}
