package controllers

import (
	"rbac/models"
	"rbac/services"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RoleController struct {
	service services.RoleService
	baseController
}

func (this *RoleController) ToView() {
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err2 := rpcService.HasRole(currentUser.Id, "root")
	if err2 != nil || !ok {
		return
	}
	tree, resTree, err := this.service.ListRoles(currentUser.Id)
	if err != nil {
		beego.Error(err)

	}

	this.Data["tree"] = tree
	this.Data["resTree"] = resTree
	this.Data["menu_index"] = "1-3"
	this.TplName = "role.tpl"
}

func (this *RoleController) Add() {
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
	o := orm.NewOrm()
	err = o.Begin()
	var role models.Role
	role.RoleName = this.GetString("role_name")
	role.Description = this.GetString("description")
	role.Pid, _ = this.GetInt("pid")
	res := this.GetStrings("res[]")
	role.GmtCreate = time.Now()
	role.GmtModified = time.Now()
	var id int64
	id, err = o.Insert(&role)
	if err != nil {
		beego.Informational(err)
		result.Msg = err.Error()
		err = o.Rollback()
	} else {
		roleRess := make([]models.RoleRes, 0, 10)
		for i := 0; i < len(res); i++ {
			resId, _ := strconv.Atoi(res[i])
			roleRes := models.RoleRes{RoleId: int(id), ResId: resId, GmtCreate: time.Now(), GmtModified: time.Now()}
			roleRess = append(roleRess, roleRes)
		}
		_, err = o.InsertMulti(100, roleRess)
		if err != nil {
			beego.Informational(err)
			result.Msg = err.Error()
			err = o.Rollback()

		} else {
			err = o.Commit()
			result.Code = 1
			result.Data = id
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RoleController) Edit() {
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
	var role models.Role
	role.Id, _ = this.GetInt("id")
	res := this.GetStrings("res[]")
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		beego.Informational(err)
		result.Msg = err.Error()
	} else {
		err = o.Read(&role)
		if err != nil {
			beego.Informational(err)
			result.Msg = err.Error()
			err = o.Rollback()
		} else {
			role.RoleName = this.GetString("role_name")
			role.Description = this.GetString("description")
			role.GmtModified = time.Now()
			_, err = o.Update(&role)
			if err != nil {
				beego.Informational(err)
				result.Msg = err.Error()
				err = o.Rollback()
			} else {
				if _, err = o.QueryTable("role_res").Filter("role_id", role.Id).Delete(); err != nil {
					beego.Informational(err)
					result.Msg = err.Error()
					err = o.Rollback()
				} else {
					roleRess := make([]models.RoleRes, 0, 10)
					for i := 0; i < len(res); i++ {
						resId, _ := strconv.Atoi(res[i])
						roleRes := models.RoleRes{RoleId: int(role.Id), ResId: resId, GmtCreate: time.Now(), GmtModified: time.Now()}
						roleRess = append(roleRess, roleRes)
					}
					if _, err = o.InsertMulti(100, roleRess); err != nil {
						beego.Informational(err)
						result.Msg = err.Error()
						err = o.Rollback()
					} else {
						err = o.Commit()
						result.Code = 1
					}
				}
			}
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RoleController) Delete() {
	result := new(models.BaseResponse)
	currentUser := this.GetSession("user").(models.User)
	rpcService := new(services.RpcService)
	ok, err := rpcService.HasRole(currentUser.Id, "root")
	if err != nil || !ok {
		result.Msg = "no permission"
		this.Data["json"] = result
		this.ServeJSON()
		return
	}
	role := new(models.Role)
	role.Id, _ = this.GetInt("id")
	o := orm.NewOrm()
	err = o.Begin()
	err = o.Read(role)
	if err != nil {
		beego.Informational(err)
		result.Msg = err.Error()
		err = o.Rollback()
	} else {
		_, err = o.Delete(role)
		if err != nil {
			beego.Informational(err)
			result.Msg = err.Error()
			err = o.Rollback()
		} else {
			roles := make([]*models.Role, 0, 10)
			qs := o.QueryTable("role").Filter("pid", role.Id)
			_, err = qs.All(&roles)
			if err != nil {
				beego.Informational(err)
				result.Msg = err.Error()
				err = o.Rollback()
			} else {
				for i := 0; i < len(roles); i++ {
					roles[i].Pid = role.Pid
					_, err = o.Update(roles[i])
					if err != nil {
						beego.Informational(err)
						result.Msg = err.Error()
						err = o.Rollback()
						break
					}
				}
				if _, err = o.QueryTable("role_res").Filter("role_id", role.Id).Delete(); err != nil {
					beego.Informational(err)
					result.Msg = err.Error()
					err = o.Rollback()
				} else {
					err = o.Commit()
					result.Code = 1
				}
			}
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RoleController) GetAllChildRoleByUserId() {
	result := new(models.BaseResponse)
	currentUser := this.GetSession("user").(models.User)
	withMe, _ := this.GetBool("withMe")
	roleIds, err := this.service.GetAllChildRoleByUserId(currentUser.Id, new(services.RpcService), withMe)
	if err != nil {
		beego.Error(err)
		result.Msg = err.Error()
	} else {
		o := orm.NewOrm()
		var roles []*models.Role
		_, err = o.QueryTable("role").Filter("id__in", roleIds).All(&roles)
		if err != nil {
			beego.Error(err)
			result.Msg = err.Error()
		} else {
			result.Code = 1
		}
	}

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *RoleController) toTree(role []*models.Role, tree *models.TreeBean) {
	for i := 0; i < len(role); i++ {
		if role[i].Pid == tree.Id {
			var childTree *models.TreeBean = &models.TreeBean{}
			childTree.Id = role[i].Id
			childTree.Label = role[i].RoleName
			childTree.Role.Id = role[i].Id
			childTree.Role.RoleName = role[i].RoleName
			childTree.Role.Description = role[i].Description
			childTree.Role.Pid = role[i].Pid
			this.toTree(role, childTree)
			if tree.Children == nil {
				tree.Children = []*models.TreeBean{}
			}
			tree.Children = append(tree.Children, childTree)
		}
	}
}
