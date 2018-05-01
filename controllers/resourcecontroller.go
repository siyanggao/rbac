package controllers

import (
	"rbac/models"
	"rbac/services"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ResourceController struct {
	baseController
}

func (this *ResourceController) ToView() {
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err := rpcService.HasRole(currentUser.Id, "root")
	if err != nil || !ok {
		return
	}
	o := orm.NewOrm()
	res := make([]*models.Resource, 100)
	qs := o.QueryTable("resource")
	_, err = qs.All(&res)
	if err != nil {
		beego.Info(err)
	}
	tree := []*models.TreeBean{}
	for i := 0; i < len(res); i++ {
		if res[i].Pid == 0 {
			childTree := &models.TreeBean{}
			childTree.Id = res[i].Id
			childTree.Label = res[i].ResName
			childTree.Res.Id = res[i].Id
			childTree.Res.ResName = res[i].ResName
			childTree.Res.ResCode = res[i].ResCode
			childTree.Res.Pid = res[i].Pid
			this.toTree(res, childTree)
			tree = append(tree, childTree)
		}
	}
	this.Data["tree"] = tree
	this.Data["menu_index"] = "1-4"
	this.TplName = "resource.tpl"
}

func (this *ResourceController) Add() {
	result := &models.BaseResponse{}
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err := rpcService.HasRole(currentUser.Id, "root")
	if err != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	resName := this.GetString("res_name")
	resCode := this.GetString("res_code")
	pId, _ := this.GetInt("pid")
	o := orm.NewOrm()
	parentRes := models.Resource{Id: pId}
	err = o.Read(&parentRes)
	if err != nil {
		result.Msg = err.Error()
	} else {
		if parentRes.ResCode == "" {
			result.Msg = "pid not exist"
		} else {
			res := models.Resource{ResName: resName, ResCode: resCode, Pid: pId, GmtCreate: time.Now(), GmtModified: time.Now()}
			id, err2 := o.Insert(&res)
			if err2 != nil {
				result.Msg = err2.Error()
			} else {
				result.Code = 1
				result.Data = id
			}
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *ResourceController) Edit() {
	result := &models.BaseResponse{}
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err := rpcService.HasRole(currentUser.Id, "root")
	if err != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	id, _ := this.GetInt("id")
	resName := this.GetString("res_name")
	resCode := this.GetString("res_code")
	o := orm.NewOrm()
	res := models.Resource{Id: id}
	err = o.Read(&res)
	if err != nil {
		result.Msg = err.Error()
	} else {
		res.ResName = resName
		res.ResCode = resCode
		res.GmtModified = time.Now()
		num, err2 := o.Update(&res)
		if err2 != nil {
			result.Msg = err2.Error()
		} else if num == 0 {
			result.Msg = "update failure"
		} else {
			result.Code = 1
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *ResourceController) Delete() {
	result := &models.BaseResponse{}
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err := rpcService.HasRole(currentUser.Id, "root")
	if err != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	res := make([]*models.Resource, 0)
	qs := o.QueryTable("resource")
	_, err = qs.All(&res)
	if err != nil {
		beego.Informational(err)
		result.Msg = err.Error()
	} else {
		ids := make([]int, 0)
		ids = append(ids, id)
		this.findResIdByPid(id, res, ids)
		var sqlWhere string
		for i := 0; i < len(ids); i++ {
			sqlWhere += strconv.Itoa(ids[i]) + ","
		}
		sqlWhere = strings.Trim(sqlWhere, "/,/")
		sql := "delete from resource where id in(" + sqlWhere + ")"
		beego.Informational(sql)
		_, errDel := o.Raw(sql).Exec()
		if errDel != nil {
			result.Msg = err.Error()
		} else {
			result.Code = 1
		}

	}

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *ResourceController) findResIdByPid(pid int, res []*models.Resource, ids []int) {
	for i := 0; i < len(res); i++ {
		if res[i].Pid == pid {
			ids = append(ids, res[i].Id)
			this.findResIdByPid(res[i].Id, res, ids)
		}
	}
}

func (this *ResourceController) toTree(res []*models.Resource, tree *models.TreeBean) {
	for i := 0; i < len(res); i++ {
		if res[i].Pid == tree.Id {
			var childTree *models.TreeBean = &models.TreeBean{}
			childTree.Id = res[i].Id
			childTree.Label = res[i].ResName
			childTree.Res.Id = res[i].Id
			childTree.Res.ResName = res[i].ResName
			childTree.Res.ResCode = res[i].ResCode
			childTree.Res.Pid = res[i].Pid
			this.toTree(res, childTree)
			if tree.Children == nil {
				tree.Children = []*models.TreeBean{}
			}
			tree.Children = append(tree.Children, childTree)
		}
	}
}
