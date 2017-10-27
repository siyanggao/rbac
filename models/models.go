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

type Role struct {
	Id          int
	RoleName    string
	Description string
	Pid         int
	GmtCreate   time.Time
	GmtModified time.Time
}

type RoleRes struct {
	Id          int
	RoleId      int
	ResId       int
	GmtCreate   time.Time
	GmtModified time.Time
}

type Depart struct {
	Id          int
	Name        string
	Pid         int
	Code        string
	Level       int
	Pids        string
	GmtCreate   time.Time
	GmtModified time.Time
}

type DepartRes struct {
	Id       int
	DepartId int
	ResId    int
	Type     int
}

type DepartRole struct {
	Id       int
	DepartId int
	RoleId   int
}

type UserDepart struct {
	Id       int
	UserId   int
	DepartId int
}

type User struct {
	Id           int
	UserName     string
	Birth        time.Time
	Gender       int
	Addr         string
	Pwd          string
	Mobile       string
	Email        string
	RealName     string
	IdentityCard string
	Status       int
	Token        string
	ClientId     string
}

type UserRole struct {
	Id     int
	UserId int
	RoleId int
	Type   int
}

type UserRes struct {
	Id     int
	UserId int
	ResId  int
	Type   int
}

type TreeBean struct {
	Id       int
	Label    string
	Disabled bool
	Res      Resource
	Role     Role
	Depart   Depart
	Children []*TreeBean
}

type BaseResponse struct {
	Code int
	Msg  string
	Data interface{}
}

type BaseModel struct {
	Id          int
	GmtCreate   time.Time
	GmtModified time.Time
}

func RegisterDB() {
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("db_username")+":"+
		beego.AppConfig.String("db_pass")+"@tcp("+
		beego.AppConfig.String("db_url")+":"+
		beego.AppConfig.String("db_port")+")/"+
		beego.AppConfig.String("db_name")+"?charset="+
		beego.AppConfig.String("db_charset"), 300)
	orm.RegisterModel(new(Resource))
	orm.RegisterModel(new(Role))
	orm.RegisterModel(new(RoleRes))
	orm.RegisterModel(new(Depart), new(DepartRes), new(DepartRole), new(UserDepart), new(User), new(UserRole), new(UserRes))
	orm.RunSyncdb("default", false, true)
}
