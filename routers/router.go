// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"zwProject/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/apply",
			beego.NSInclude(
				&controllers.ApplyController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/load_record",
			beego.NSInclude(
				&controllers.LoadRecordController{},
			),
		),
		beego.NSNamespace("/inventory",
			beego.NSInclude(
				&controllers.InventoryController{},
			),
		),
	)
	beego.AddNamespace(ns)
	beego.Router("/wxaccesstoken", &controllers.WXAccessTokenController{})
	beego.Router("/wx", &controllers.WXController{})
	beego.Router("/wxmini/login", &controllers.WXMiniLoginController{})
	beego.Router("/wxmini/register", &controllers.WXMiniRegisterController{})
	beego.Router("/wxmini/load", &controllers.WXMiniLoadController{})
	beego.Router("/wxmini/unload", &controllers.WXMiniUnLoadController{})
	beego.Router("/wxmini/pickeritem", &controllers.WXMiniPickerItemController{})
	beego.Router("/wxmini/filter_load_record", &controllers.WXMiniFilterLoadRecordController{})
}
