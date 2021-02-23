package controllers

import (
	// "encoding/json"
	"fmt"
	"encoding/base64"
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

	// public_key_n := "26096591567557089821307549233404438514323316819882733790173025852589043479063638473371782409819473121611942939328823991764302201326430603263196226482799058768293218147913198410924615192343752596206585922931854347874350566610717382730269410379415778608784086728946198901983793060718637627435153248390779225276402766629855817216750438527691727285639781421387552721542115193873966108252849077727102685328813998626798612584316734647889335414611736669324640667395930780983553839590893679335782669128582296760709055242002130082808750074278277519296253097467405357102081348934283812070731572293049764026118304233627606580771"
	// public_key_e := "65537"
	// certificate := "wo shi hao ren"

	ReturnData := make(map[string]interface{})

	ReturnData["pk_n"] = public_key_n
	ReturnData["pk_e"] = public_key_e
	ReturnData["certificate"] = certificate

	c.Data["json"] = ReturnData
	c.ServeJSON() //响应前端
	c.StopRun()

}

func (c *SessionKeyController) SkPost() {
	ReturnData := make(map[string]interface{})

	date := c.GetString("encrypted_session_key")

	encrypted_session_key,_ := base64.StdEncoding.DecodeString(date)

	fmt.Println(encrypted_session_key)
	// fmt.Println([]byte(encrypted_session_key))

	get_session_key("1", string(encrypted_session_key[:]))

	ReturnData["result"] = true

	c.Data["json"] = ReturnData
	c.ServeJSON() //响应前端
	c.StopRun()
}