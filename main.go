package main

import (
	_ "rbac/routers"

	"rbac/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	orm.Debug = true
	models.RegisterDB()
	beego.Run()

}
