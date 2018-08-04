package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"zwProject/common"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"zwProject/models"
	"zwProject/db"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type WXMiniLoginController struct {
	beego.Controller
}
type WXMiniRegisterController struct {
	beego.Controller
}
func (c *WXMiniLoginController) Get() {
	code:=c.GetString("code")
	fmt.Println("code is :"+code)
	if code == ""{
		fmt.Println("get code err")
		return
	}
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=wxfe7815dd10b97a64&secret=242a3a1bbb058c3d95efcd14445dccac&js_code="+code+"&grant_type=authorization_code")
	if err != nil {
		fmt.Println("get openid err")
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}
	var respObj LoginResponseBody
	json.Unmarshal(body,&respObj)

	beego.BeeLogger.Info("openid:"+respObj.Openid)

	//判断openid是否在user库里，如果不在显示申请角色，在的话显示菜单
	user:=models.User{Openid:respObj.Openid}
	err = db.GetOrm().Read(&user,"openid")
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
		//TODO 申请角色，写入角色申请表
		//c.Ctx.WriteString("haha")
		c.Data["json"] = &CodeBody{"0"}
		c.ServeJSON()
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	} else {
		fmt.Println(user)
		fmt.Println("查询到用户"+strconv.Itoa(user.Id))


		//更新数据库中昵称，头像地址
		nickName:=c.GetString("nickName")
		avatarUrl:=c.GetString("avatarUrl")
		user.Username = nickName
		user.AvatarUrl = avatarUrl

		if num, err := db.GetOrm().Update(&user); err == nil {
			fmt.Println("更新昵和地址,影响行数:"+strconv.Itoa(int(num)))
		}
		//TODO 显示菜单
	}
}
func (c *WXMiniLoginController) Post() {
	var codebody CodeBody
	common.ProcJsonRequest(c.Ctx.ResponseWriter.ResponseWriter,c.Ctx.Request,&codebody)
	fmt.Print("code:"+codebody.Code)


}

func (c *WXMiniRegisterController) Post() {
	var codebody CodeBody
	common.ProcJsonRequest(c.Ctx.ResponseWriter.ResponseWriter,c.Ctx.Request,&codebody)
	fmt.Print("code:"+codebody.Code)


}


type CodeBody struct {
	Code string
}
type LoginResponseBody struct {
	Session_key string
	Openid string
}


