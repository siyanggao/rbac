package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"rbac/models"
	"rbac/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserService struct {
	baseService
}

func (this *UserService) GetUserByDepart(departId int, currentUser models.User) ([]models.User, error) {
	o := orm.NewOrm()
	ok, err := new(DepartService).IsMyChildDepart(currentUser.Id, []int{departId}, true)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if !ok {
		return nil, errors.New("no permission")
	}
	var user []models.User
	_, err = o.Raw("select t2.* from user_depart t left join user t2 on t.user_id=t2.id where t.depart_id=?", departId).QueryRows(&user)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return user, nil
}

func (this *UserService) ListUser(page int, length int, user models.User, currentUser models.User, rpcService *RpcService) ([]models.User, error) {
	start := new(baseService).ParsePage(page, length)
	o := orm.NewOrm()
	var users []models.User
	qs := o.QueryTable("user")
	if len(user.UserName) > 0 {
		qs = qs.Filter("user_name", user.UserName)
	}
	child, err := this.GetChildByUser(currentUser.Id, rpcService)
	if err != nil {
		return nil, err
	}

	qs = qs.Filter("id__in", child)

	_, err = qs.Limit(length, start).All(&users)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return users, nil
}

func (this *UserService) Add(u models.User) (int, error) {
	o := orm.NewOrm()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte("123456"))
	u.Pwd = hex.EncodeToString(md5Ctx.Sum(nil))
	o.Begin()
	id, err := o.Insert(&u)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return 0, err
	}
	if len(u.Avatar) > 0 {
		path, filename, subfix := utils.ParsePathFile(u.Avatar)
		path = strings.Replace(path, "/tmp", "/avatar", 1)
		filename = strconv.Itoa(int(id))

		if err2 := os.Rename(u.Avatar, path+filename+subfix); err2 != nil {
			o.Rollback()
			return 0, err2
		}
		u.Avatar = path + filename + subfix
		_, err = o.Update(&u, "Avatar")
		if err != nil {
			beego.Error(err)
			o.Rollback()
			return 0, err
		}
	}
	o.Commit()
	return int(id), nil
}

func (this *UserService) Edit(u models.User, currentUser models.User, rpcService *RpcService) error {
	ok, err := this.IsMyChildUser(currentUser.Id, []int{u.Id}, rpcService)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	o := orm.NewOrm()
	if strings.Contains(u.Avatar, "/tmp") {
		path, filename, subfix := utils.ParsePathFile(u.Avatar)
		path = strings.Replace(path, "/tmp", "/avatar", 1)
		filename = strconv.Itoa(u.Id)
		if err2 := os.Rename(u.Avatar, path+filename+subfix); err2 != nil {
			return err2
		}
		u.Avatar = path + filename + subfix
	}
	if _, err := o.Update(&u, "UserName", "Mobile", "Gender", "RealName", "IdentityCard", "Avatar"); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func (this *UserService) Delete(id int, currentUser models.User, rpcService *RpcService) error {

	ok, err := this.IsMyChildUser(currentUser.Id, []int{id}, rpcService)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	o := orm.NewOrm()
	u := models.User{Id: id}
	if _, err := o.Delete(&u); err != nil {
		beego.Error(err)
		return err
	}
	bm.Delete("GetChildByUser_" + strconv.Itoa(currentUser.Id))
	bm.Delete("GetDepartByUser_" + strconv.Itoa(currentUser.Id))
	bm.Delete("GetRoleByUser_" + strconv.Itoa(currentUser.Id))
	bm.Delete("GetResByUser_" + strconv.Itoa(currentUser.Id))
	return nil
}

func (this *UserService) AllotDepart(userId int, departId []int, currentUser models.User, rpcService *RpcService) error {
	bm.Delete("GetChildByUser" + strconv.Itoa(userId))
	ok, err := this.IsMyChildUser(currentUser.Id, []int{userId}, rpcService)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	ok, err = new(DepartService).IsMyChildDepart(currentUser.Id, departId, false)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("not my child,no permission")
	}
	var depart []*models.Depart
	o := orm.NewOrm()
	_, err = o.QueryTable("depart").Filter("id__in", departId).All(&depart)
	if err != nil {
		beego.Error(err)
		return err
	}
	o.Begin()
	_, err = o.QueryTable("user_depart").Filter("user_id__in", userId).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	var userDepart = []models.UserDepart{}
	for i := 0; i < len(depart); i++ {
		userDepart = append(userDepart, models.UserDepart{UserId: userId, DepartId: depart[i].Id})
	}
	_, err = o.InsertMulti(100, userDepart)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	bm.Delete("GetDepartByUser_" + strconv.Itoa(currentUser.Id))
	return nil
}

func (this *UserService) AllotRole(userId int, roleId []int, currentUser models.User, rpcService *RpcService) error {
	ok, err := this.IsMyChildUser(currentUser.Id, []int{userId}, rpcService)
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
	o := orm.NewOrm()
	var departRole []models.Role
	_, err = o.Raw("select t4.* from user t "+
		"left join user_depart t2 on t.id=t2.user_id "+
		"left join depart_role t3 on t2.depart_id=t3.depart_id "+
		"left join role t4 on t3.role_id=t4.id "+
		"where t.id=?", userId).QueryRows(&departRole)
	if err != nil {
		beego.Error(err)
		return err
	}
	var unableRole []*models.Role = make([]*models.Role, 0, 10) //unable role where exist in depart's role
	for _, item := range departRole {
		var isInDepartRole bool //roleId is or not exist in depart's role
		for index, item2 := range roleId {
			if item.Id == item2 {
				roleId = append(roleId[:index], roleId[index+1:]...) //remove role where exist in user's depart
				isInDepartRole = true
			}
		}
		if !isInDepartRole {
			unableRole = append(unableRole, &item)
		}
	}

	//struct userRole data
	var userRole []*models.UserRole = make([]*models.UserRole, 0, 10)
	for _, item := range roleId {
		userRole = append(userRole, &models.UserRole{UserId: userId, RoleId: item, Type: 0})
	}
	for _, item := range unableRole {
		userRole = append(userRole, &models.UserRole{UserId: userId, RoleId: item.Id, Type: 1})
	}
	o.Begin()
	_, err = o.QueryTable("user_role").Filter("user_id", userId).Delete()
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return err
	}
	if len(userRole) > 0 {
		_, err = o.InsertMulti(100, userRole)
		if err != nil {
			beego.Error(err)
			o.Rollback()
			return err
		}
	}

	o.Commit()
	bm.Delete("GetRoleByUser_" + strconv.Itoa(currentUser.Id))
	return nil
}

func (this *UserService) AllotRes(userId int, resId []int, currentUser models.User, rpcService *RpcService) error {
	ok, err := this.IsMyChildUser(currentUser.Id, []int{userId}, rpcService)
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
	beego.Informational(myRes)
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
	var res []*models.Resource
	depart, err := rpcService.GetDepartByUser(userId)
	if err != nil {
		beego.Error(err)
		return err
	}
	var unableDepartRes []*models.DepartRes = make([]*models.DepartRes, 0, 10)
	for _, item := range depart {
		tmpRes, err2 := rpcService.GetResByDepart(item.Id)
		if err2 != nil {
			beego.Error(err2)
			return err2
		}
		res = append(res, tmpRes...)
		var tmpUnableDepartRes []*models.DepartRes
		_, err2 = o.QueryTable("depart_res").Filter("depart_id", item.Id).Filter("type", 1).All(&tmpUnableDepartRes)
		if err2 != nil {
			beego.Error(err2)
			return err2
		}
		unableDepartRes = append(unableDepartRes, tmpUnableDepartRes...)
	}
	role, err2 := rpcService.GetRoleByUser(userId)
	var roleRes []*models.Resource = make([]*models.Resource, 0, 10)
	if err2 != nil {
		beego.Error(err2)
		return err2
	}
	for _, item := range role {
		tmpRes, err3 := rpcService.GetResByRole(item.Id)
		if err3 != nil {
			beego.Error(err3)
			return err3
		}
		roleRes = append(roleRes, tmpRes...)
	}
	for index, item := range roleRes {
		for _, item2 := range unableDepartRes {
			if item.Id == item2.ResId {
				if index == len(roleRes)-1 {
					roleRes = append(roleRes[0:index])
				} else {
					roleRes = append(roleRes[0:index], roleRes[index+1:]...)
				}
			}
		}
	}
	res = append(res, roleRes...)
	for i := 0; i < len(res); i++ {
		for j := i; j < len(res); j++ {
			if res[i].Id > res[j].Id {
				tmp := res[i]
				res[i] = res[j]
				res[j] = tmp
			}
		}
	}
	dstRes := make([]*models.Resource, 0, 10)
	if len(res) > 0 {
		dstRes = append(dstRes, res[0])
		for _, item := range res {
			if dstRes[len(dstRes)-1].Id != item.Id {
				dstRes = append(dstRes, item)
			}
		}
	}
	beego.Informational(dstRes)

	var unableRes []*models.Resource = make([]*models.Resource, 0, 10) //unable res where not exist in depart,role's res
	for _, item := range dstRes {
		var isInDepartRoleRes bool //resId is or not exist in depart,role's resource
		for index, item2 := range resId {
			if item.Id == item2 {
				if index == len(resId)-1 { //remove res where exist in depart,role's res
					resId = append(resId[0:index])
				} else {
					resId = append(resId[0:index], resId[index+1:]...)
				}
				beego.Informational(index, resId)
				isInDepartRoleRes = true
				index--
			}
		}
		if !isInDepartRoleRes {
			beego.Informational(item)
			unableRes = append(unableRes, item)
		}
	}
	//struct departRes data
	var userRes []*models.UserRes = make([]*models.UserRes, 0, 10)
	for _, item := range resId {
		userRes = append(userRes, &models.UserRes{UserId: userId, ResId: item, Type: 0})
	}
	for _, item := range unableRes {
		userRes = append(userRes, &models.UserRes{UserId: userId, ResId: item.Id, Type: 1})
	}

	o.Begin()
	_, err = o.QueryTable("user_res").Filter("user_id", userId).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	if len(userRes) > 0 {
		_, err = o.InsertMulti(100, userRes)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	bm.Delete("GetResByUser_" + strconv.Itoa(currentUser.Id))
	return nil
}

func (this *UserService) Count(user models.User) (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	if len(user.UserName) > 0 {
		qs = qs.Filter("user_name", user.UserName)
	}
	count, err := qs.Count()
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	return int(count), nil
}

func (this *UserService) GetByUsername(username string) (user models.User, err error) {
	o := orm.NewOrm()
	user.UserName = username
	err = o.Read(&user, "UserName")
	if err != nil {
		beego.Error(err)
	}
	return
}

func (this *UserService) GetChildByUser(userId int, rpcService *RpcService) (users []int, err error) {
	if bm.IsExist("GetChildByUser_" + strconv.Itoa(userId)) {
		return bm.Get("GetChildByUser_" + strconv.Itoa(userId)).([]int), nil
	}
	o := orm.NewOrm()
	depart := make([]*models.Depart, 0, 10)
	_, err = o.QueryTable("depart").All(&depart)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	var myDepartIds *[]int = new([]int)
	myDepart, err3 := rpcService.GetDepartByUser(userId)
	if err3 != nil {
		return nil, err3
	}
	departService := new(DepartService)
	for _, item := range myDepart {
		*myDepartIds = append(*myDepartIds, item.Id)
		departService.findDepartIdByPid(item.Id, depart, myDepartIds)
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

	var userDepart []*models.UserDepart
	_, err = o.QueryTable("user_depart").Filter("depart_id__in", dstMyDepartIds).All(&userDepart)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	var noDepartUser []*models.User
	_, err = o.Raw("select t.* from user t left join user_depart t2 on t.id=t2.user_id where t2.id is null").QueryRows(&noDepartUser)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	users = make([]int, 0, 10)
	for _, item := range userDepart {
		users = append(users, item.UserId)
	}
	for _, item := range noDepartUser {
		users = append(users, item.Id)
	}
	bm.Put("GetChildByUser_"+strconv.Itoa(userId), users, 3600*24*7*time.Second)
	return users, nil
}

func (this *UserService) IsMyChildUser(userId int, user []int, rpcService *RpcService) (bool, error) {
	myChild, err := this.GetChildByUser(userId, rpcService)
	if err != nil {
		return false, err
	}
	for _, item := range user {
		var isMyChild bool
		for _, item2 := range myChild {
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
