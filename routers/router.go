// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"myproject/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/permisson",
			beego.NSInclude(
				&controllers.PermissonController{},
			),
		),

		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RoleController{},
			),
		),

		beego.NSNamespace("/role_permisson",
			beego.NSInclude(
				&controllers.RolePermissonController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),

		beego.NSNamespace("/user_role",
			beego.NSInclude(
				&controllers.UserRoleController{},
			),
		),
	)
	beego.AddNamespace(ns)
	beego.Router("/wxaccesstoken", &controllers.WXAccessTokenController{})
	beego.Router("/wx", &controllers.WXController{})
}
