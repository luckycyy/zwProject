package db

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

var o orm.Ormer

func InitDB(){
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass:=beego.AppConfig.String("mysqlpass")
	mysqlserver := beego.AppConfig.String("mysqlserver")
	mysqlport := beego.AppConfig.String("mysqlport")
	mysqldb := beego.AppConfig.String("mysqldb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqluser+":"+mysqlpass+"@tcp("+mysqlserver+":"+mysqlport+")/"+mysqldb+"?loc=Asia%2FShanghai")
	orm.Debug = true
	o =orm.NewOrm()
	o.Using("default")
}

func GetOrm() orm.Ormer{
	return o
}