package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"zwProject/common"
	"io/ioutil"
	"net/http"
)

type WXMiniLoginController struct {
	beego.Controller
}
func (c *WXMiniLoginController) Get() {
	code:=c.GetString("code")
	fmt.Println("code1 is :"+code)
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=wxfe7815dd10b97a64&secret=242a3a1bbb058c3d95efcd14445dccac&js_code="+code+"&grant_type=authorization_code")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}

	fmt.Println("resp:"+string(body))
}
func (c *WXMiniLoginController) Post() {
	var codebody CodeBody
	common.ProcJsonRequest(c.Ctx.ResponseWriter.ResponseWriter,c.Ctx.Request,&codebody)
	fmt.Print("code:"+codebody.Code)


}

type CodeBody struct {
	Code string
}