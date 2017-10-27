package controllers

import (
	"rbac/models"
	"rbac/services"

	"github.com/astaxie/beego"
)

type RpcController struct {
	service services.RpcService
	baseController
}

func (this *RpcController) HasRes(userId int, resCode ...string) {
	result := new(models.BaseResponse)
	ok, err := this.service.HasRes(userId, resCode...)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	result.Data = ok
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RpcController) HasRole(userId int, roleName ...string) {
	result := new(models.BaseResponse)
	ok, err := this.service.HasRole(userId, roleName...)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
	}
	result.Data = ok
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RpcController) GetResByRole() {
	roleId, _ := this.GetInt("id")
	result := models.BaseResponse{}
	res, err := this.service.GetResByRole(roleId)
	if err != nil {
		beego.Informational(err)
		result.Msg = err.Error()

	} else {
		result.Code = 1
		result.Data = res
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RpcController) GetRoleByDepart() {
	departId, _ := this.GetInt("id")
	result := models.BaseResponse{}
	role, err := this.service.GetRoleByDepart(departId)
	if err != nil {
		beego.Informational(err)
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = role
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RpcController) GetResByDepart() {
	departId, _ := this.GetInt("id")
	result := new(models.BaseResponse)
	res, err := this.service.GetResByDepart(departId)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = res
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RpcController) GetDepartByUser() {
	userId, _ := this.GetInt("Id")
	result := new(models.BaseResponse)
	depart, err := this.service.GetDepartByUser(userId)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = depart
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RpcController) GetRoleByUser() {
	userId, _ := this.GetInt("Id")
	result := new(models.BaseResponse)
	role, err := this.service.GetRoleByUser(userId)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = role
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RpcController) GetResByUser() {
	userId, _ := this.GetInt("Id")
	result := new(models.BaseResponse)
	res, err := this.service.GetResByUser(userId)
	if err != nil {
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = res
	}
	this.Data["json"] = result
	this.ServeJSON()
}
