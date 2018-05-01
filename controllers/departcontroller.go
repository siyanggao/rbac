package controllers

import (
	"rbac/models"
	"rbac/services"
	"strconv"

	"github.com/astaxie/beego"
)

type DepartController struct {
	service services.DepartService
	baseController
}

//func (this *DepartController) GetDepartByRole() {
//	roleId, _ := this.GetInt("roleId")
//	depart, err := this.service.GetDepartByRole(roleId)
//	result := new(models.BaseResponse)
//	if err != nil {
//		result.Msg = err.Error()
//	} else {
//		result.Code = 1
//		result.Data = depart
//	}
//	this.Data["json"] = result
//	this.ServeJSON()
//}

func (this *DepartController) ToView() {
	rpcService := new(services.RpcService)
	user := this.GetSession("user").(models.User)
	ok, err2 := rpcService.HasRes(user.Id, "depart:detail")
	if err2 != nil || !ok {
		this.TplName = "login.tpl"
		return
	}
	tree, roleTree, resTree, err := this.service.ListDeparts(user.Id, rpcService)
	if err != nil {
		beego.Error(err)
	}

	this.Data["tree"] = tree
	this.Data["roleTree"] = roleTree
	this.Data["resTree"] = resTree
	this.Data["menu_index"] = "1-2"
	this.TplName = "depart.tpl"
}

func (this *DepartController) Add() {
	result := &models.BaseResponse{}
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err2 := rpcService.HasRes(currentUser.Id, "depart:add")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}

	var depart *models.Depart = new(models.Depart)
	depart.Name = this.GetString("name")
	depart.Pid, _ = this.GetInt("pid")
	id, err := this.service.Add(depart, currentUser)
	if err != nil {
		beego.Error(err)
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = id
	}

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *DepartController) Edit() {
	result := &models.BaseResponse{}
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err2 := rpcService.HasRes(currentUser.Id, "depart:edit")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}

	depart := new(models.Depart)
	depart.Id, _ = this.GetInt("id")
	depart.Name = this.GetString("name")
	err := this.service.Edit(depart, currentUser)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *DepartController) Delete() {
	result := new(models.BaseResponse)
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err2 := rpcService.HasRes(currentUser.Id, "depart:del")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}

	depart := new(models.Depart)
	depart.Id, _ = this.GetInt("id")
	err := this.service.Delete(depart, currentUser, rpcService)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *DepartController) AllotRole() {
	result := new(models.BaseResponse)
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err2 := rpcService.HasRes(currentUser.Id, "depart:allotrole")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	departId, _ := this.GetInt("id")
	role := this.GetStrings("role[]")

	roleInt := make([]int, 0, 10)
	for i := 0; i < len(role); i++ {
		value, _ := strconv.Atoi(role[i])
		roleInt = append(roleInt, value)
	}
	err := this.service.AllotRole(departId, roleInt, currentUser, rpcService)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *DepartController) AllotRes() {
	result := new(models.BaseResponse)
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err2 := rpcService.HasRes(currentUser.Id, "depart:allotres")
	if err2 != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	departId, _ := this.GetInt("id")
	res := this.GetStrings("res[]")
	resInt := make([]int, 0, 10)
	for i := 0; i < len(res); i++ {
		value, _ := strconv.Atoi(res[i])
		resInt = append(resInt, value)
	}

	err := this.service.AllotRes(departId, resInt, currentUser, rpcService)
	if err != nil {
		beego.Error(err)
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *DepartController) GetParents() {
	result := new(models.BaseResponse)
	currentUser := this.GetSession("user").(models.User)
	parents, err := this.service.GetParents(currentUser.Id)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = parents
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *DepartController) GetByRoleName() {
	result := new(models.BaseResponse)
	currentUser := this.GetSession("user").(models.User)
	roleName := this.GetString("roleName")
	departs, err := this.service.GetByRoleName(currentUser.Id, roleName)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = departs
	}
	this.Data["json"] = result
	this.ServeJSON()
}
