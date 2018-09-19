package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"zwProject/common"
	"zwProject/db"
	"zwProject/models"
)

type WXMiniLoginController struct {
	beego.Controller
}
type WXMiniRegisterController struct {
	beego.Controller
}
type WXMiniLoadController struct {
	beego.Controller
}
type WXMiniUnLoadController struct {
	beego.Controller
}
type WXMiniPickerItemController struct {
	beego.Controller
}
type WXMiniFilterLoadRecordController struct {
	beego.Controller
}

func (c *WXMiniFilterLoadRecordController) Get() {
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	strEndDate := endDate
	describe := c.GetString("describe")
	driver := c.GetString("driver")
	station := c.GetString("station")
	productName := c.GetString("productName")

	t, _ := time.Parse("2006-01-02", endDate)
	t = t.AddDate(0, 0, 1) //end日期 加一天，防止选不出当天的
	endDate = t.Format("2006-01-02")

	//var lists []orm.ParamsList

	var loadRecord []*models.LoadRecord

	qs := db.GetOrm().QueryTable("load_record")
	if ("" != startDate) && ("0000-00-00" != strEndDate) {
		qs = qs.Filter("create_time__gte", startDate)
		fmt.Println("s")
	}
	if ("" != endDate) && ("0000-00-00" != strEndDate) {
		qs = qs.Filter("create_time__lte", endDate)
		fmt.Println("e")
	}

	if "" != describe {
		qs = qs.Filter("describe", describe)
	}
	if "" != driver {
		qs = qs.Filter("creator", driver)
	}
	if "" != station {
		qs = qs.Filter("station", station)
	}
	if "" != productName {
		qs = qs.Filter("product_name", productName)
	}

	num, err := qs.OrderBy("-create_time").All(&loadRecord)
	fmt.Printf("Returned Rows Num: %s, %s", num, err)
	/*
	num, err := qs.ValuesList(&lists)
	if err == nil {
		fmt.Printf("Result Nums: %d\n", num)
		for _, row := range lists {
			fmt.Println(row)
		}
	}*/

	c.Data["json"] = loadRecord
	c.ServeJSON()

}
func (c *WXMiniLoginController) Get() {
	code := c.GetString("code")
	fmt.Println("code is :" + code)
	if code == "" {
		fmt.Println("get code err")
		return
	}
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=wxfe7815dd10b97a64&secret=242a3a1bbb058c3d95efcd14445dccac&js_code=" + code + "&grant_type=authorization_code")
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
	json.Unmarshal(body, &respObj)

	beego.BeeLogger.Info("openid:" + respObj.Openid)

	//判断openid是否在user库里，如果不在显示申请角色，在的话显示菜单
	//user:=models.User{Openid:respObj.Openid}
	//err = db.GetOrm().Read(&user,"openid")
	user := models.User{}
	err = db.GetOrm().QueryTable("user").Filter("openid", respObj.Openid).One(&user)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		fmt.Printf("返回多个user不是一个")
	}

	if err == orm.ErrNoRows {
		fmt.Println("user中查询不到")

		apply := models.Apply{Openid: respObj.Openid}
		err = db.GetOrm().Read(&apply, "openid")

		var role string
		if err == orm.ErrNoRows {
			role = "unregister"
		} else {
			role = "待审核"
		}

		c.Data["json"] = &LoginResult{respObj.Openid, role, "", "",""}
		c.ServeJSON()
		return
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
		return
	} else {
		fmt.Println(user)
		fmt.Println("查询到用户" + strconv.Itoa(user.Id))

		//更新数据库中昵称，头像地址
		nickName := c.GetString("nickName")
		avatarUrl := c.GetString("avatarUrl")
		user.Nickname = nickName
		user.AvatarUrl = avatarUrl

		if num, err := db.GetOrm().Update(&user, "nickname", "avatar_url"); err == nil {
			fmt.Println("更新昵和地址,影响行数:" + strconv.Itoa(int(num)))
		}
		//查询picker条目
		pickerItemsJsonStr, err := json.Marshal(getPickerItems())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(pickerItemsJsonStr)
		c.Data["json"] = &LoginResult{respObj.Openid, user.Role, user.Username, user.Describe,string(pickerItemsJsonStr)}
		c.ServeJSON()
	}
}
func (c *WXMiniLoginController) Post() {
	var codebody CodeBody
	common.ProcJsonRequest(c.Ctx.ResponseWriter.ResponseWriter, c.Ctx.Request, &codebody)
	fmt.Print("cod1e:" + codebody.Code)

}

func (c *WXMiniRegisterController) Post() {
	var apply models.Apply
	//err:=json.Unmarshal(c.Ctx.Input.RequestBody, &apply)
	err := c.ParseForm(&apply)
	if err != nil {
		fmt.Println("parseForm err,", err)
		return
	}
	apply.CreateTime = time.Now().UTC()
	applyId, err := models.AddApply(&apply)
	fmt.Println("apply：", apply)
	fmt.Println("applyId：", applyId)
	if err != nil {
		fmt.Println("addApply err,", err)
		return
	}

	c.Data["json"] = &CodeBody{"1"}
	c.ServeJSON()

}

func (c *WXMiniLoadController) Post() {
	var loadRecord models.LoadRecord
	//err:=json.Unmarshal(c.Ctx.Input.RequestBody, &apply)
	err := c.ParseForm(&loadRecord)
	if err != nil {
		fmt.Println("parseForm err,", err)
		return
	}
	loadRecord.CreateTime = time.Now().UTC()
	loadRecordId, err := models.AddLoadRecord(&loadRecord)
	fmt.Println("loadRecord：", loadRecord)
	fmt.Println("loadRecordId：", loadRecordId)
	if err != nil {
		fmt.Println("addloadRecord err,", err)
		return
	}

	//扣减库存
	res, err := db.GetOrm().Raw("UPDATE inventory SET num = num - ? WHERE product_name = ? AND station = ?", loadRecord.Num, loadRecord.ProductName, loadRecord.Station).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		fmt.Println("update minus inventory affected nums: ", num)
		if num == 0{
			c.Data["json"] = &CodeBody{"2"}//站点与产品不匹配
		}else{
			c.Data["json"] = &CodeBody{"1"}//成功
		}
	}else{
		c.Data["json"] = &CodeBody{"err"}
	}
	//TODO 修改为原生语句 或者存入整型
	//num, err := db.GetOrm().QueryTable("inventory").Filter("product_name",loadRecord.ProductName).Filter("station",loadRecord.Station).Update(orm.Params{
	//	"num": orm.ColValue(orm.ColMinus, loadRecord.Num),
	//})
	//fmt.Println("影响num：", num)

	c.ServeJSON()

}
func (c *WXMiniUnLoadController) Put() {

	record := models.LoadRecord{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &record)

	if err != nil {
		fmt.Println("Unmarshal err,", err)
		return
	}

	fmt.Println(record)
	if record.Location2 != "" {
		fmt.Println("jr location2")
		if num, err := db.GetOrm().Update(&record, "location2"); err == nil {
			fmt.Println("recordID:" + strconv.Itoa(record.Id) + "更新卸车位置" + record.Location2 + ",影响行数:" + strconv.Itoa(int(num)))
		}
	} else {
		fmt.Println("jr else")
		if num, err := db.GetOrm().Update(&record, "is_unload"); err == nil {
			fmt.Println("recordID:" + strconv.Itoa(record.Id) + "更新卸车,影响行数:" + strconv.Itoa(int(num)))
		}
	}

	c.Data["json"] = &CodeBody{"OK"}
	c.ServeJSON()
}

func (c *WXMiniPickerItemController) Get() {
	value := c.GetString("key")
	fmt.Println("key is :" + value)
	if value == "" {
		fmt.Println("get key err")
		return
	}

	rs := getPickerItems()
	c.Data["json"] = rs
	c.ServeJSON()

}

func getPickerItems() []interface{} {
	var products orm.ParamsList
	var stations orm.ParamsList
	num, err := db.GetOrm().Raw("SELECT '请选择产品名称' UNION SELECT DISTINCT product_name FROM inventory").ValuesFlat(&products)
	if err == nil && num > 0 {
		fmt.Println(products) // []{"1","2","3",...}
	}
	num, err = db.GetOrm().Raw("SELECT '请选择装车站点' UNION SELECT DISTINCT station FROM inventory").ValuesFlat(&stations)
	if err == nil && num > 0 {
		fmt.Println(stations) // []{"1","2","3",...}
	}
	rs := make([]interface{}, 2)
	rs[0] = products
	rs[1] = stations
	return rs
}

type CodeBody struct {
	Code string
}
type LoginResult struct {
	Openid      string
	Role        string
	Username    string
	Describe    string
	PickerItems string
}
type LoginResponseBody struct {
	Session_key string
	Openid      string
}
