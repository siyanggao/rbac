package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"rbac/models"
	"rbac/services"
	"strconv"

	"github.com/astaxie/beego"
)

type UserController struct {
	service services.UserService
	baseController
}

func (this *UserController) GetUserByDepart() {
	currentUser := this.GetSession("user").(models.User)
	departId, _ := this.GetInt("departId")
	depart, err := this.service.GetUserByDepart(departId, currentUser)
	result := new(models.BaseResponse)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = depart
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) ToView() {
	rpcService := new(services.RpcService)
	currentUser := this.GetSession("user").(models.User)
	//	ok, err2 := rpcService.HasRes(currentUser.Id, "user:detail")
	//	if err2 != nil || !ok {

	//		this.TplName = "login.tpl"
	//		return
	//	}
	users, err := this.service.ListUser(1, 10, models.User{}, currentUser, rpcService)
	if err != nil {
		beego.Error(err)
	}
	var count int
	count, err = this.service.Count(models.User{})
	depart, role, res, err2 := new(services.DepartService).ListDeparts(currentUser.Id, rpcService)
	if err2 != nil {

	}
	this.Data["tableData"] = users
	this.Data["totalSize"] = count
	this.Data["menu_index"] = "1-1"
	this.Data["depart"] = depart
	this.Data["role"] = role
	this.Data["res"] = res

	this.TplName = "user.tpl"
}

func (this *UserController) Page() {
	rpcService := new(services.RpcService)
	currentUser := this.GetSession("user").(models.User)
	//	ok, err2 := rpcService.HasRes(currentUser.Id, "user:detail")
	//	if err2 != nil || !ok {

	//		this.TplName = "login.tpl"
	//		return
	//	}
	username := this.GetString("search_username")
	user := models.User{UserName: username}
	page, _ := this.GetInt("currentPage")
	pageSize, _ := this.GetInt("pageSize")
	totalSize, _ := this.GetInt("totalSize")
	result := new(models.BaseResponse)
	users, err := this.service.ListUser(page, pageSize, user, currentUser, rpcService)
	if err != nil {
		result.Msg = err.Error()
	} else {
		var count int
		if totalSize == 0 {
			count, err = this.service.Count(user)
		}
		result.Code = 1
		result.Data = struct {
			Users []models.User
			Count int
		}{
			Users: users,
			Count: count,
		}
	}
	this.Data["json"] = result
	this.ServeJSON()

}

func (this *UserController) Login() {
	username := this.GetString("username")
	pwd := this.GetString("pwd")
	user, err := this.service.GetByUsername(username)
	result := new(models.BaseResponse)
	if err != nil {
		result.Msg = err.Error()
	} else {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(pwd))
		pwd = hex.EncodeToString(md5Ctx.Sum(nil))
		if pwd == user.Pwd {
			this.SetSession("user", user)
			result.Code = 1
			result.Data = user
		} else {
			result.Msg = "password error"
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) LoginView() {
	this.TplName = "login.tpl"
}

func (this *UserController) Add() {
	result := new(models.BaseResponse)
	ok, err2 := new(services.RpcService).HasRes(this.GetSession("user").(models.User).Id, "user:add")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	u := models.User{}
	if err := this.ParseForm(&u); err != nil {
		beego.Error(err)
		result.Msg = err.Error()
	} else {
		var id int
		id, err = this.service.Add(u)
		if err != nil {
			result.Msg = err.Error()
		} else {
			result.Code = 1
			result.Data = id
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) Edit() {
	result := new(models.BaseResponse)
	rpcService := new(services.RpcService)
	currentUser := this.GetSession("user").(models.User)
	ok, err2 := rpcService.HasRes(currentUser.Id, "user:edit")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	u := models.User{}
	if err := this.ParseForm(&u); err != nil {
		beego.Error(err)
		result.Msg = err.Error()
	} else {
		if err = this.service.Edit(u, currentUser, rpcService); err != nil {
			result.Msg = err.Error()
		} else {
			result.Code = 1
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) Delete() {
	result := new(models.BaseResponse)
	rpcService := new(services.RpcService)
	currentUser := this.GetSession("user").(models.User)
	ok, err2 := rpcService.HasRes(currentUser.Id, "user:del")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	id, _ := this.GetInt("id")

	if err := this.service.Delete(id, currentUser, rpcService); err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) AllotDepart() {
	result := new(models.BaseResponse)
	rpcService := new(services.RpcService)
	currentUser := this.GetSession("user").(models.User)
	ok, err2 := rpcService.HasRes(currentUser.Id, "user:allotdepart")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	userId, _ := this.GetInt("Id")
	depart := this.GetStrings("depart[]")
	departInt := make([]int, 0, 10)
	for i := 0; i < len(depart); i++ {
		value, _ := strconv.Atoi(depart[i])
		departInt = append(departInt, value)
	}
	err := this.service.AllotDepart(userId, departInt, currentUser, rpcService)

	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) AllotRole() {
	result := new(models.BaseResponse)
	rpcService := new(services.RpcService)
	currentUser := this.GetSession("user").(models.User)
	ok, err2 := rpcService.HasRes(currentUser.Id, "user:allotrole")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	userId, _ := this.GetInt("Id")
	role := this.GetStrings("role[]")

	roleInt := make([]int, 0, 10)
	for i := 0; i < len(role); i++ {
		value, _ := strconv.Atoi(role[i])
		roleInt = append(roleInt, value)
	}

	err := this.service.AllotRole(userId, roleInt, currentUser, rpcService)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UserController) AllotRes() {
	result := new(models.BaseResponse)
	rpcService := new(services.RpcService)
	currentUser := this.GetSession("user").(models.User)
	ok, err2 := rpcService.HasRes(currentUser.Id, "user:allotres")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	userId, _ := this.GetInt("Id")
	res := this.GetStrings("res[]")
	resInt := make([]int, 0, 10)
	for i := 0; i < len(res); i++ {
		value, _ := strconv.Atoi(res[i])
		resInt = append(resInt, value)
	}

	err := this.service.AllotRes(userId, resInt, currentUser, rpcService)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	this.Data["json"] = result
	this.ServeJSON()
}
