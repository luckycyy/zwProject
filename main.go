package main

import (
	_ "zwProject/routers"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"zwProject/db"
)

func init() {
	db.InitDB()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

