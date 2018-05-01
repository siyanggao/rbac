package services

import (
	"rbac/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RoleService struct {
	baseService
}

func (this *RoleService) ListRoles(userId int) ([]*models.TreeBean, []*models.TreeBean, error) {
	o := orm.NewOrm()
	role := make([]*models.Role, 100)
	qs := o.QueryTable("role")
	_, err := qs.All(&role)
	if err != nil {
		beego.Error(err.Error())
		return nil, nil, err
	}
	rpcService := new(RpcService)
	myRole, err4 := rpcService.GetRoleByUser(userId)
	if err4 != nil {
		return nil, nil, err4
	}
	var myRoleIds *[]int = new([]int)
	for _, item := range myRole {
		*myRoleIds = append(*myRoleIds, item.Id)
		this.FindRoleIdByPid(item.Id, role, myRoleIds)
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
	tree := []*models.TreeBean{}
	for i := 0; i < len(role); i++ {
		beego.Informational(role[i])
		if role[i].Pid == 0 {
			childTree := &models.TreeBean{}
			childTree.Id = role[i].Id
			childTree.Label = role[i].RoleName
			childTree.Disabled = true
			childTree.Role.Id = role[i].Id
			childTree.Role.RoleName = role[i].RoleName
			childTree.Role.Description = role[i].Description
			childTree.Role.Pid = role[i].Pid
			this.toTree(role, childTree, dstMyRoleIds)
			tree = append(tree, childTree)
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
		return tree, nil, err2
	}
	myRes, err5 := rpcService.GetResByUser(userId)
	if err5 != nil {
		return nil, nil, err5
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
	return tree, resTree, nil
}

func (this *RoleService) IsMyChildRole(userId int, roleId []int, withMe bool) (bool, error) {
	myRoleIds, err := this.GetAllChildRoleByUserId(userId, new(RpcService), withMe)
	if err != nil {
		beego.Error(err)
		return false, err
	}
	for _, item := range myRoleIds {
		for _, item2 := range roleId {
			if item == item2 {
				return true, nil
			}
		}
	}
	return false, nil
}

func (this *RoleService) GetAllChildRoleByUserId(userId int, rpcService *RpcService, withMe bool) ([]int, error) {
	if bm.IsExist("GetAllChildRoleByUserId" + strconv.Itoa(userId)) {
		return bm.Get("GetAllChildRoleByUserId" + strconv.Itoa(userId)).([]int), nil
	}
	myRole, err := rpcService.GetRoleByUser(userId)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	var role []*models.Role
	_, err = o.QueryTable("role").All(&role)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	var myRoleIds *[]int = new([]int)
	for _, item := range myRole {
		if withMe {
			*myRoleIds = append(*myRoleIds, item.Id)
		}
		this.FindRoleIdByPid(item.Id, role, myRoleIds)
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
	bm.Put("GetAllChildRoleByUserId"+strconv.Itoa(userId), dstMyRoleIds, 3600*24*7*time.Second)
	return dstMyRoleIds, nil
}

func (this *RoleService) toTree(role []*models.Role, tree *models.TreeBean, dstMyRoleIds []int) {
	for i := 0; i < len(role); i++ {
		if role[i].Pid == tree.Id {
			var childTree *models.TreeBean = &models.TreeBean{}
			childTree.Id = role[i].Id
			childTree.Label = role[i].RoleName
			childTree.Disabled = true
			childTree.Role.Id = role[i].Id
			childTree.Role.RoleName = role[i].RoleName
			childTree.Role.Description = role[i].Description
			childTree.Role.Pid = role[i].Pid
			this.toTree(role, childTree, dstMyRoleIds)
			if tree.Children == nil {
				tree.Children = []*models.TreeBean{}
			}
			tree.Children = append(tree.Children, childTree)
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
}

func (this *RoleService) FindRoleIdByPid(pid int, role []*models.Role, ids *[]int) {
	for i := 0; i < len(role); i++ {
		if role[i].Pid == pid {
			*ids = append(*ids, role[i].Id)
			this.FindRoleIdByPid(role[i].Id, role, ids)
		}
	}
}
