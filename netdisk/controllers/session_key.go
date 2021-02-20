package controllers

import (
	// "encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

type SessionKeyController struct {
	beego.Controller
}

func (c *SessionKeyController) SkGet() {

	// 从sgx获取证书

	public_key_n, public_key_e, certificate := server_hello()

	fmt.Println("pbk_n:{}",public_key_n)
	fmt.Println("pbk_e:{}",public_key_e)
	fmt.Println(certificate)

	// public_key_n := []byte{179,97,236,123,127,76,215,155,180,156,152,19,185,88,37,255,255,240,104,156,8,210,10,227,54,58,167,156,30,232,115,17,180,138,65,53,203,143,188,160,17,175,191,7,241,238,63,117,64,137,241,16,48,237,164,157,254,225,126,123,192,217,225,18,96,197,114,155,168,88,144,177,195,46,77,155,222,231,252,219,40,112,170,230,32,102,106,167,92,49,36,137,163,207,0,100,250,224,230,227,131,186,86,244,191,251,206,39,64,23,47,152,149,137,215,220,73,119,104,209,183,96,80,132,63,182,240,191,181,96,87,188,143,103,254,99,197,211,246,195,97,67,69,7,192,118,218,115,125,244,155,65,156,84,57,242,35,5,12,195,231,138,193,137,202,23,23,163,30,179,44,88,103,229,136,62,226,86,129,254,139,242,197,190,100,54,28,18,212,248,160,118,170,181,193,159,178,12,145,211,72,30,73,79,142,13,8,105,71,205,254,96,2,177,16,16,94,17,119,199,106,16,88,19,247,95,208,151,149,77,15,166,98,247,204,113,32,202,132,248,50,240,33,97,170,104,219,197,203,255,37,228,69,252,169,195}
	// public_key_e := []byte{1,0,1}
	// certificate := "wo shi hao ren"

	// fmt.Println(len(string(public_key_n)))
	// fmt.Println(len(public_key_n))
	// fmt.Println(string(public_key_n))
	// fmt.Println(string(public_key_e))

	// fmt.Println(public_key_e)

	ReturnData := make(map[string]interface{})

	ReturnData["pk_n"] = []byte(public_key_n)
	ReturnData["pk_e"] = []byte(public_key_e)
	ReturnData["certificate"] = certificate

	c.Data["json"] = ReturnData
	c.ServeJSON() //响应前端
	c.StopRun()

}

func (c *SessionKeyController) SkPost() {
	c.TplName = "personalcenter.html"
}
