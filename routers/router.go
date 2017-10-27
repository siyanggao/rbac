package routers

import (
	"rbac/controllers"
	"rbac/models"
	"rbac/services"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.ResourceController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.RpcController{})
	beego.AutoRouter(&controllers.DepartController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.InsertFilter("/*", beego.BeforeRouter, FilterLogin)
	beego.InsertFilter("/resource/*", beego.BeforeRouter, FilterRoot)
	beego.InsertFilter("/role/*", beego.BeforeRouter, FilterRoot)
}

var FilterLogin = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("user").(models.User)
	if !ok && ctx.Request.RequestURI != "/user/loginview" && ctx.Request.RequestURI != "/user/login" {
		ctx.Redirect(302, "/user/loginview")
	}
}

var FilterRoot = func(ctx *context.Context) {
	user := ctx.Input.Session("user")
	has, _ := new(services.RpcService).HasRole(user.(models.User).Id, "root")
	if !has {
		ctx.Redirect(302, "/user/loginview")
	}

}
