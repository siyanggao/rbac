package services

import (
	"errors"
	"rbac/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DepartService struct {
	baseService
}

func (this *DepartService) GetDepartByRole(roleId int) ([]models.Depart, error) {
	o := orm.NewOrm()
	var depart []models.Depart
	_, err := o.Raw("select t2.* from depart_role t left join depart t2 on t.depart_id=t2.id where t.role_id=?", roleId).QueryRows(&depart)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return depart, nil
}

func (this *DepartService) ListDeparts(userId int, rpcService *RpcService) ([]*models.TreeBean, []*models.TreeBean, []*models.TreeBean, error) {
	o := orm.NewOrm()
	depart := make([]*models.Depart, 0, 10)
	qs := o.QueryTable("depart")
	_, err := qs.All(&depart)
	if err != nil {
		beego.Error(err.Error())
		return nil, nil, nil, err
	}

	var myDepartIds *[]int = new([]int)
	myDepart, err3 := rpcService.GetDepartByUser(userId)
	if err3 != nil {
		return nil, nil, nil, err3
	}
	for _, item := range myDepart {
		*myDepartIds = append(*myDepartIds, item.Id)
		this.findDepartIdByPid(item.Id, depart, myDepartIds)
	}
	for i := 0; i < len(*myDepartIds); i++ {
		for j := i + 1; j < len(*myDepartIds); j++ {
			if (*myDepartIds)[i] > (*myDepartIds)[j] {
				tmp := (*myDepartIds)[i]
				(*myDepartIds)[i] = (*myDepartIds)[j]
				(*myDepartIds)[j] = tmp
			}
		}
	}
	var dstMyDepartIds []int
	if len(*myDepartIds) > 0 {
		dstMyDepartIds = []int{(*myDepartIds)[0]}
	}
	for _, item := range *myDepartIds {
		if item != dstMyDepartIds[len(dstMyDepartIds)-1] {
			dstMyDepartIds = append(dstMyDepartIds, item)
		}
	}

	tree := []*models.TreeBean{}
	for i := 0; i < len(depart); i++ {
		if depart[i].Pid == 0 {
			childTree := &models.TreeBean{}
			childTree.Id = depart[i].Id
			childTree.Label = depart[i].Name
			childTree.Disabled = true
			childTree.Depart.Id = depart[i].Id
			childTree.Depart.Name = depart[i].Name
			childTree.Depart.Pid = depart[i].Pid
			this.toTree(depart, childTree, dstMyDepartIds)
			tree = append(tree, childTree)
			if len(dstMyDepartIds) != 0 {
				for _, item := range dstMyDepartIds {
					if childTree.Depart.Id == item {
						childTree.Disabled = false
						break
					}
				}
			}
		}
	}

	role := make([]*models.Role, 0, 10)
	_, err = o.QueryTable("role").All(&role)
	if err != nil {
		beego.Error(err)
		return tree, nil, nil, err
	}
	myRole, err4 := rpcService.GetRoleByUser(userId)
	if err4 != nil {
		return nil, nil, nil, err4
	}
	var myRoleIds *[]int = new([]int)
	for _, item := range myRole {
		*myRoleIds = append(*myRoleIds, item.Id)
		new(RoleService).FindRoleIdByPid(item.Id, role, myRoleIds)
	}

	for i := 0; i < len(*myRoleIds); i++ {
		for j := i + 1; j < len(*myRoleIds); j++ {
			if (*myRoleIds)[i] > (*myRoleIds)[j] {
				tmp := (*myRoleIds)[i]
				(*myRoleIds)[i] = (*myRoleIds)[j]
				(*myRoleIds)[j] = tmp
			}
		}
	}
	var dstMyRoleIds []int
	if len(*myRoleIds) > 0 {
		dstMyRoleIds = []int{(*myRoleIds)[0]}
	}
	for _, item := range *myRoleIds {
		if item != dstMyRoleIds[len(dstMyRoleIds)-1] {
			dstMyRoleIds = append(dstMyRoleIds, item)
		}
	}
	roleTree := []*models.TreeBean{}
	for i := 0; i < len(role); i++ {
		if role[i].Pid == 0 {
			childTree := &models.TreeBean{}
			childTree.Id = role[i].Id
			childTree.Label = role[i].RoleName
			childTree.Disabled = true
			childTree.Role.Id = role[i].Id
			childTree.Role.RoleName = role[i].RoleName
			childTree.Role.Description = role[i].Description
			childTree.Role.Pid = role[i].Pid
			new(RoleService).toTree(role, childTree, dstMyRoleIds)
			roleTree = append(roleTree, childTree)
			if len(dstMyRoleIds) != 0 {
				for _, item := range dstMyRoleIds {
					if childTree.Role.Id == item {
						childTree.Disabled = false
						break
					}
				}
			}
		}
	}

	res := make([]*models.Resource, 100)
	_, err2 := o.QueryTable("resource").All(&res)
	if err2 != nil {
		beego.Error(err2)
		return tree, nil, nil, err2
	}
	myRes, err5 := rpcService.GetResByUser(userId)
	if err5 != nil {
		return nil, nil, nil, err5
	}

	resTree := []*models.TreeBean{}
	for i := 0; i < len(res); i++ {
		if res[i].Pid == 0 {
			childTree := &models.TreeBean{}
			childTree.Id = res[i].Id
			childTree.Label = res[i].ResName
			childTree.Disabled = true
			childTree.Res.Id = res[i].Id
			childTree.Res.ResName = res[i].ResName
			childTree.Res.ResCode = res[i].ResCode
			childTree.Res.Pid = res[i].Pid
			new(ResourceService).toTree(res, childTree, myRes)
			resTree = append(resTree, childTree)
			if len(myRes) != 0 {
				for _, item := range myRes {
					if childTree.Depart.Id == item.Id {
						childTree.Disabled = false
						break
					}
				}
			}
		}
	}
	return tree, roleTree, resTree, nil
}

func (this *DepartService) Add(depart *models.Depart, currentUser models.User) (int, error) {
	myDepartIds, err := this.GetAllChildDepartByUserId(currentUser.Id, true)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	var isMyChild bool
	for _, item := range *myDepartIds {
		if item == depart.Pid {
			isMyChild = true
			break
		}
	}
	if !isMyChild {
		return 0, errors.New("not mychild,no permission")
	}
	o := orm.NewOrm()
	parent := models.Depart{Id: depart.Pid}
	err = o.Read(&parent)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	depart.Level = parent.Level + 1
	if len(parent.Pids) == 0 {
		depart.Pids = strconv.Itoa(parent.Id)
	} else {
		depart.Pids = parent.Pids + "," + strconv.Itoa(parent.Id)
	}

	var child []*models.Depart
	child, err = this.getDirectChild(&parent, o)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	var max int
	max, err = this.getMaxNum(child)
	if max+1 < 10 {
		depart.Code = "00" + strconv.Itoa(max+1)
	} else if max+1 < 100 {
		depart.Code = "0" + strconv.Itoa(max+1)
	} else {
		depart.Code = "" + strconv.Itoa(max+1)
	}
	var id int64
	id, err = o.Insert(depart)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (this *DepartService) Edit(depart *models.Depart, currentUser models.User) error {
	myDepartIds, err := this.GetAllChildDepartByUserId(currentUser.Id, true)
	if err != nil {
		beego.Error(err)
		return err
	}
	var isMyChild bool
	for _, item := range *myDepartIds {
		if item == depart.Pid {
			isMyChild = true
			break
		}
	}
	if !isMyChild {
		return errors.New("not mychild,no permission")
	}
	o := orm.NewOrm()
	_, err = o.Update(depart, "Name")
	if err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func (this *DepartService) Delete(depart *models.Depart, currentUser models.User, rpcService *RpcService) error {
	ok, err := this.IsMyChildDepart(currentUser.Id, []int{depart.Pid}, false)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	o := orm.NewOrm()
	o.Begin()
	var departs []*models.Depart
	_, err = o.QueryTable("depart").All(&departs)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return err
	}
	var ids = make([]int, 0, 10)
	this.findDepartIdByPid(depart.Id, departs, &ids)
	ids = append(ids, depart.Id)
	_, err = o.QueryTable("depart").Filter("id__in", ids).Delete()
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return err
	}
	_, err = o.QueryTable("depart_res").Filter("depart_id__in", ids).Delete()
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return err
	}
	_, err = o.QueryTable("depart_role").Filter("depart_id__in", ids).Delete()
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return err
	}
	_, err = o.QueryTable("user_depart").Filter("depart_id__in", ids).Delete()
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return err
	}
	err = o.Commit()
	bm.Delete("GetAllChildDepartByUserId_" + strconv.Itoa(currentUser.Id))
	bm.Delete("GetRoleByDepart_" + strconv.Itoa(depart.Id))
	bm.Delete("GetResByDepart_" + strconv.Itoa(depart.Id))
	return nil
}

func (this *DepartService) AllotRole(departId int, roleId []int, currentUser models.User, rpcService *RpcService) error {
	ok, err := this.IsMyChildDepart(currentUser.Id, []int{departId}, false)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	ok, err = new(RoleService).IsMyChildRole(currentUser.Id, roleId, true)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	var role []*models.Role
	o := orm.NewOrm()
	_, err = o.QueryTable("role").Filter("id__in", roleId).All(&role)
	if err != nil {
		beego.Error(err)
		return err
	}
	o.Begin()
	_, err = o.QueryTable("depart_role").Filter("depart_id__in", roleId).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	var departRole = []models.DepartRole{}
	for i := 0; i < len(role); i++ {
		departRole = append(departRole, models.DepartRole{DepartId: departId, RoleId: role[i].Id})
	}
	_, err = o.InsertMulti(100, departRole)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	bm.Delete("GetRoleByDepart_" + strconv.Itoa(departId))
	return nil
}

func (this *DepartService) AllotRes(departId int, resId []int, currentUser models.User, rpcService *RpcService) error {
	ok, err := this.IsMyChildDepart(currentUser.Id, []int{departId}, false)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	myRes, err2 := rpcService.GetResByUser(currentUser.Id)
	if err2 != nil {
		return err
	}
	for _, item := range resId {
		var isMyChild bool
		for _, item2 := range myRes {
			if item == item2.Id {
				isMyChild = true
				break
			}
		}
		if !isMyChild {
			return errors.New("not my child,no permission")
		}
	}
	o := orm.NewOrm()
	var roleRes []models.Resource
	_, err = o.Raw("select t4.* from depart t "+
		"left join depart_role t2 on t.id=t2.depart_id "+
		"left join role_res t3 on t2.role_id=t3.role_id "+
		"left join resource t4 on t3.res_id=t4.id "+
		"where t.id=?", departId).QueryRows(&roleRes)
	if err != nil {
		return err
	}
	var unableRes []*models.Resource = make([]*models.Resource, 0, 10) //unable res where exist in role's res
	for _, item := range roleRes {
		var isInRoleRes bool //resId is or not exist in role's resource
		for index, item2 := range resId {
			if item.Id == item2 {
				resId = append(resId[:index], resId[index+1:]...) //remove res where exist in role's res
				isInRoleRes = true
			}
		}
		if !isInRoleRes {
			unableRes = append(unableRes, &item)
		}
	}

	//struct departRes data
	var departRes []*models.DepartRes = make([]*models.DepartRes, 0, 10)
	for _, item := range resId {
		departRes = append(departRes, &models.DepartRes{DepartId: departId, ResId: item, Type: 0})
	}
	for _, item := range unableRes {
		departRes = append(departRes, &models.DepartRes{DepartId: departId, ResId: item.Id, Type: 1})
	}
	o.Begin()
	_, err = o.QueryTable("depart_res").Filter("depart_id", departId).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	if len(departRes) > 0 {
		_, err = o.InsertMulti(100, departRes)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	bm.Delete("GetResByDepart_" + strconv.Itoa(departId))
	return nil
}

func (this *DepartService) IsMyChildDepart(userId int, departId []int, withMe bool) (bool, error) {
	myDepartIds, err := this.GetAllChildDepartByUserId(userId, withMe)
	if err != nil {
		beego.Error(err)
		return false, err
	}
	for _, item := range departId {
		var isMyChild bool
		for _, item2 := range *myDepartIds {
			if item == item2 {
				isMyChild = true
				break
			}
		}
		if !isMyChild {
			return false, nil
		}
	}
	return true, nil
}

func (this *DepartService) getDirectChild(parent *models.Depart, o orm.Ormer) ([]*models.Depart, error) {
	child := make([]*models.Depart, 0, 10)
	_, err := o.QueryTable("depart").Filter("pid", parent.Id).All(&child)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return child, nil
}

func (this *DepartService) GetAllChildDepartByUserId(userId int, withMe bool) (*[]int, error) {
	if bm.IsExist("GetAllChildDepartByUserId" + strconv.Itoa(userId)) {
		return bm.Get("GetAllChildDepartByUserId" + strconv.Itoa(userId)).(*[]int), nil
	}
	rpcService := new(RpcService)
	var myDepartIds *[]int = new([]int)
	myDepart, err := rpcService.GetDepartByUser(userId)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	var depart []*models.Depart
	_, err = o.QueryTable("depart").All(&depart)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	for _, item := range myDepart {
		if withMe {
			*myDepartIds = append(*myDepartIds, item.Id)
		}
		this.findDepartIdByPid(item.Id, depart, myDepartIds)
	}
	for i := 0; i < len(*myDepartIds); i++ {
		for j := i + 1; j < len(*myDepartIds); j++ {
			if (*myDepartIds)[i] > (*myDepartIds)[j] {
				tmp := (*myDepartIds)[i]
				(*myDepartIds)[i] = (*myDepartIds)[j]
				(*myDepartIds)[j] = tmp
			}
		}
	}
	var dstMyDepartIds []int
	if len(*myDepartIds) > 0 {
		dstMyDepartIds = []int{(*myDepartIds)[0]}
	}
	for _, item := range *myDepartIds {
		if item != dstMyDepartIds[len(dstMyDepartIds)-1] {
			dstMyDepartIds = append(dstMyDepartIds, item)
		}
	}
	bm.Put("GetAllChildDepartByUserId"+strconv.Itoa(userId), myDepartIds, 3600*24*7*time.Second)
	return myDepartIds, nil
}

func (this *DepartService) getMaxNum(depart []*models.Depart) (int, error) {
	var max int
	var numStrs []int = make([]int, 0, 10)
	for i := 0; i < len(depart); i++ {
		tmp := strings.Split(depart[i].Code, "_")
		tmpNum, err := strconv.Atoi(strings.TrimLeft(tmp[len(tmp)-1], "0"))
		if err != nil {
			return max, err
		}
		numStrs = append(numStrs, tmpNum)
		if max < numStrs[i] {
			max = numStrs[i]
		}
	}
	return max, nil
}

func (this *DepartService) toTree(depart []*models.Depart, tree *models.TreeBean, myDepart []int) {
	for i := 0; i < len(depart); i++ {
		if depart[i].Pid == tree.Id {
			var childTree *models.TreeBean = &models.TreeBean{}
			childTree.Id = depart[i].Id
			childTree.Label = depart[i].Name
			childTree.Disabled = true
			childTree.Depart.Id = depart[i].Id
			childTree.Depart.Name = depart[i].Name
			childTree.Depart.Pid = depart[i].Pid
			this.toTree(depart, childTree, myDepart)
			if tree.Children == nil {
				tree.Children = []*models.TreeBean{}
			}
			tree.Children = append(tree.Children, childTree)
			if len(myDepart) != 0 {
				for _, item := range myDepart {
					if childTree.Depart.Id == item {
						childTree.Disabled = false
						break
					}
				}
			}
		}
	}

}

func (this *DepartService) findDepartIdByPid(pid int, depart []*models.Depart, ids *[]int) {
	for i := 0; i < len(depart); i++ {
		if depart[i].Pid == pid {
			*ids = append(*ids, depart[i].Id)
			this.findDepartIdByPid(depart[i].Id, depart, ids)
		}
	}
}
