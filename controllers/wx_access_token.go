package controllers

import (
	"github.com/astaxie/beego"
	"zwProject/wx"
)

type WXAccessTokenController struct {
	beego.Controller
}

func (c *WXAccessTokenController) Get() {
	wx.GetWXAccessToken()
	c.Ctx.WriteString(wx.GetWXAccessToken())
}