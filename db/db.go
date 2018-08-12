package db

import (
	"github.com/astaxie/beego/orm"
)

var o orm.Ormer

func InitDB(){
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(101.201.71.156:3306)/myproject?loc=Asia%2FShanghai")
	orm.Debug = true
	o =orm.NewOrm()
	o.Using("default")
}

func GetOrm() orm.Ormer{
	return o
}