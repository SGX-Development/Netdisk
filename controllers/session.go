package controllers

import (

)
//这个结构体有问题，后面修
type UserStatus struct {
	userName string
	userid int
	islogin bool
}

func Islogin(status interface{}) bool {
	return !(status == nil || (status != nil && !status.(UserStatus).islogin))
}

func UserName(status interface{}) string {
	return status.(UserStatus).userName
}