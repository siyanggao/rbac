package controllers

import (
	"rbac/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ResourceController struct {
	baseController
}

func (this *ResourceController) ToView() {
	o := orm.NewOrm()
	res := make([]*models.Resource, 0)
	qs := o.QueryTable("resource")
	_, err := qs.All(&res)
	if err != nil {
		beego.Info(err)
	}

	this.Data["res"] = res
	this.TplName = "resource.tpl"
}

func toTree(res []*models.Resource, tree TreeBean) {
	if tree == nil {
		for i := 0; i < len(res); i++ {
			if res.Pid == 0 {
				tree.Id = res.Id
				tree.label = res.ResName
				toTree(res, tree)
				break
			}
		}
	} else {

		for i := 0; i < len(res); i++ {
			if res.Pid == tree.Id {
				childTree
			}
		}
	}

}
