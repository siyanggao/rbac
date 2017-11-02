package services

import (
	"rbac/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RpcService struct {
	baseService
}

func (this *RpcService) HasRes(userId int, resCode ...string) (bool, error) {

	ok, err := this.HasRole(userId, "root")
	if err == nil && ok {
		return true, nil
	}
	res, err2 := this.GetResByUser(userId)
	if err2 != nil {
		return false, err2
	}
	for _, item := range res {
		for _, item2 := range resCode {
			if item.ResCode == item2 {
				return true, nil
			}
		}
	}
	return false, nil
}

func (this *RpcService) HasRole(userId int, roleName ...string) (bool, error) {
	role, err := this.GetRoleByUser(userId)
	if err != nil {
		return false, err
	}
	for _, item := range role {
		for _, item2 := range roleName {
			if item.RoleName == item2 || item.RoleName == "root" {
				return true, nil
			}
		}
	}
	return false, nil
}

func (this *RpcService) GetResByRole(roleId int) ([]*models.Resource, error) {
	if bm.IsExist("GetResByRole_" + strconv.Itoa(roleId)) {
		return bm.Get("GetResByRole_" + strconv.Itoa(roleId)).([]*models.Resource), nil
	}
	var res []*models.Resource
	o := orm.NewOrm()
	role := models.Role{Id: roleId}
	err := o.Read(&role)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if role.RoleName == "root" {
		_, err = o.QueryTable("resource").All(&res)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		return res, nil
	}

	_, err = o.Raw("select t2.* from role_res t left join resource t2 on t.res_id=t2.id where t.role_id=?", roleId).QueryRows(&res)
	if err != nil {
		beego.Informational(err)
		return nil, err

	} else {
		bm.Put("GetResByRole_"+strconv.Itoa(roleId), res, 3600*24*7*time.Second)
		return res, nil
	}
}

func (this *RpcService) GetRoleByDepart(departId int) ([]models.Role, error) {
	if bm.IsExist("GetRoleByDepart_" + strconv.Itoa(departId)) {
		return bm.Get("GetRoleByDepart_" + strconv.Itoa(departId)).([]models.Role), nil
	}
	o := orm.NewOrm()
	var role []models.Role
	_, err := o.Raw("select t.* from role t left join depart_role t2 on t.id=t2.role_id where t2.depart_id=?", departId).QueryRows(&role)
	bm.Put("GetRoleByDepart_"+strconv.Itoa(departId), role, 3600*24*7*time.Second)
	return role, err
}

func (this *RpcService) GetResByDepart(departId int) ([]*models.Resource, error) {
	if bm.IsExist("GetResByDepart_" + strconv.Itoa(departId)) {
		return bm.Get("GetResByDepart_" + strconv.Itoa(departId)).([]*models.Resource), nil
	}
	o := orm.NewOrm()
	//select from depart's role's res
	var roleRes []*models.Resource
	_, err := o.Raw("select t4.* from depart t "+
		"left join depart_role t2 on t.id=t2.depart_id "+
		"left join role_res t3 on t2.role_id=t3.role_id "+
		"left join resource t4 on t3.res_id=t4.id "+
		"where t.id=?", departId).QueryRows(&roleRes)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	//remove unable from roleres
	var departRes []*models.DepartRes
	_, err = o.QueryTable("depart_res").Filter("depart_id", departId).Filter("type", 1).All(&departRes)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	for _, item := range departRes {
		for index, item2 := range roleRes {
			if item2.Id == item.ResId {
				roleRes = append(roleRes[:index], roleRes[index+1:]...)
			}
		}
	}
	//select from depart's res
	var res []*models.Resource
	_, err = o.Raw("select t2.* from depart_res t "+
		"left join resource t2 on t.res_id=t2.id "+
		"where t.type=0 and t.depart_id=?", departId).QueryRows(&res)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	//data
	res = append(res, roleRes...)
	//sort,de-weight
	for i := 0; i < len(res); i++ {
		for j := i; j < len(res); j++ {
			if res[i].Id > res[j].Id {
				tmp := res[i]
				res[i] = res[j]
				res[j] = tmp
			}
		}
	}
	var dstRes []*models.Resource
	if len(res) > 0 {
		dstRes = []*models.Resource{res[0]}
		for _, item := range res {
			if item != dstRes[len(dstRes)-1] {
				dstRes = append(dstRes, item)
			}
		}
	}
	bm.Put("GetResByDepart_"+strconv.Itoa(departId), dstRes, 3600*24*7*time.Second)
	return dstRes, err
}

func (this *RpcService) GetDepartByUser(userId int) (depart []models.Depart, err error) {
	if bm.IsExist("GetDepartByUser_" + strconv.Itoa(userId)) {
		return bm.Get("GetDepartByUser_" + strconv.Itoa(userId)).([]models.Depart), nil
	}
	o := orm.NewOrm()
	_, err = o.Raw("select t.* from depart t left join user_depart t2 on t.id=t2.depart_id where t2.user_id=?", userId).QueryRows(&depart)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	bm.Put("GetDepartByUser_"+strconv.Itoa(userId), depart, 3600*24*7*time.Second)
	return depart, nil
}

func (this *RpcService) GetRoleByUser(userId int) ([]*models.Role, error) {
	if bm.IsExist("GetRoleByUser_" + strconv.Itoa(userId)) {
		return bm.Get("GetRoleByUser_" + strconv.Itoa(userId)).([]*models.Role), nil
	}
	o := orm.NewOrm()
	//select from user's depart's role
	var departRole []*models.Role
	_, err := o.Raw("select t4.* from user t "+
		"left join user_depart t2 on t.id=t2.user_id "+
		"left join depart_role t3 on t2.depart_id=t3.depart_id "+
		"left join role t4 on t3.role_id=t4.id "+
		"where t.id=?", userId).QueryRows(&departRole)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	//remove unable from departRoles
	var userRole []*models.UserRole
	_, err = o.QueryTable("user_role").Filter("user_id", userId).Filter("type", 1).All(&userRole)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	for _, item := range userRole {
		for index, item2 := range departRole {
			if item2.Id == item.RoleId {
				departRole = append(departRole[:index], departRole[index+1:]...)
			}
		}
	}
	//select from user's role
	var role []*models.Role
	_, err = o.Raw("select t2.* from user_role t "+
		"left join role t2 on t.role_id=t2.id "+
		"where t.type=0 and t.user_id=?", userId).QueryRows(&role)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	//data
	role = append(role, departRole...)
	//sort,de-weight
	for i := 0; i < len(role); i++ {
		for j := i; j < len(role); j++ {
			if role[i].Id > role[j].Id {
				tmp := role[i]
				role[i] = role[j]
				role[j] = tmp
			}
		}
	}
	var dstRole []*models.Role
	if len(role) > 0 {
		dstRole = []*models.Role{role[0]}
		for _, item := range role {
			if item != dstRole[len(dstRole)-1] {
				dstRole = append(dstRole, item)
			}
		}
	}
	bm.Put("GetRoleByUser_"+strconv.Itoa(userId), dstRole, 3600*24*7*time.Second)
	return dstRole, err
}

func (this *RpcService) GetResByUser(userId int) (res []*models.Resource, err error) {
	if bm.IsExist("GetResByUser_" + strconv.Itoa(userId)) {
		return bm.Get("GetResByUser_" + strconv.Itoa(userId)).([]*models.Resource), nil
	}
	o := orm.NewOrm()
	depart, err := this.GetDepartByUser(userId)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	var unableDepartRes []*models.DepartRes = make([]*models.DepartRes, 0, 10)
	for _, item := range depart {
		tmpRes, err2 := this.GetResByDepart(item.Id)
		if err2 != nil {
			beego.Error(err2)
			return nil, err2
		}
		res = append(res, tmpRes...)
		var tmpUnableDepartRes []*models.DepartRes
		_, err2 = o.QueryTable("depart_res").Filter("depart_id", item.Id).Filter("type", 1).All(&tmpUnableDepartRes)
		if err2 != nil {
			beego.Error(err2)
			return nil, err2
		}
		unableDepartRes = append(unableDepartRes, tmpUnableDepartRes...)
	}

	role, err2 := this.GetRoleByUser(userId)
	if err2 != nil {
		beego.Error(err2)
		return nil, err2
	}
	for _, item := range role {
		tmpRes, err3 := this.GetResByRole(item.Id)
		if err3 != nil {
			beego.Error(err3)
			return nil, err3
		}
		res = append(res, tmpRes...)
	}
	//remove unable
	for index, item := range res {
		for _, item2 := range unableDepartRes {
			if item.Id == item2.ResId {
				//res = append(res[:index], res[index+1:]...)
				if index == len(res)-1 {
					res = append(res[0:index])
				} else {
					res = append(res[0:index], res[index+1:]...)
				}
			}
		}
	}

	var userRes []*models.Resource
	_, err = o.Raw("select t.* from resource t left join user_res t2 on t.id=t2.res_id where type=0 and t2.user_id=?", userId).QueryRows(&userRes)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	res = append(res, userRes...)

	var unableRes []*models.UserRes
	_, err = o.QueryTable("user_res").Filter("user_id", userId).Filter("type", 1).All(&unableRes)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	for index, item := range res {
		for _, item2 := range unableRes {
			if item.Id == item2.ResId {
				//res = append(res[:index], res[index+1:]...)
				if index == len(res)-1 {
					res = append(res[0:index])
				} else {
					res = append(res[0:index], res[index+1:]...)
				}
			}
		}
	}
	//sort,do-weight
	for i := 0; i < len(res); i++ {
		for j := i; j < len(res); j++ {
			if res[i].Id > res[j].Id {
				tmp := res[i]
				res[i] = res[j]
				res[j] = tmp
			}
		}
	}
	var dstRes []*models.Resource
	if len(res) > 0 {
		dstRes = []*models.Resource{res[0]}
		for _, item := range res {
			if item != dstRes[len(dstRes)-1] {
				dstRes = append(dstRes, item)
			}
		}
	}
	bm.Put("GetResByUser_"+strconv.Itoa(userId), dstRes, 3600*24*7*time.Second)
	return dstRes, nil
}
