package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Resource struct {
	Id      int
	ResName string
	ResCode string
	Pid     int

	GmtCreate   time.Time
	GmtModified time.Time
}

type TreeBean struct {
	Id       int
	Label    string
	Children []*TreeBean
}

func RegisterDB() {
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("db_username")+":"+
		beego.AppConfig.String("db_pass")+"@tcp("+
		beego.AppConfig.String("db_url")+":"+
		beego.AppConfig.String("db_port")+")/"+
		beego.AppConfig.String("db_name")+"?charset="+
		beego.AppConfig.String("db_charset"), 300)
	orm.RegisterModel(new(Resource))
	orm.RunSyncdb("default", false, true)
}
