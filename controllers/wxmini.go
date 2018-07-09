package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/astaxie/beego/context"
	"net/http"
)

type WXMiniController struct {
	beego.Controller
}

func (c *WXMiniController) Get() {
	checkSignature(c.Ctx.ResponseWriter,c.Ctx.Request)
}

func (c *WXMiniController) Post() {
	fmt.Println("2222")
	procCodeRequest(c.Ctx.ResponseWriter,c.Ctx.Request)

}
func procCodeRequest(response *context.Response, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var codebody CodeBody
	json.Unmarshal(body, &codebody)
	fmt.Print("code:"+codebody.Code)
}
type CodeBody struct {
	Code string
}