package services

import (
	"rbac/models"
)

type ResourceService struct {
	baseService
}

func (this *ResourceService) toTree(res []*models.Resource, tree *models.TreeBean, myRes []*models.Resource) {
	for i := 0; i < len(res); i++ {
		if res[i].Pid == tree.Id {
			var childTree *models.TreeBean = &models.TreeBean{}
			childTree.Id = res[i].Id
			childTree.Label = res[i].ResName
			childTree.Disabled = true
			childTree.Res.Id = res[i].Id
			childTree.Res.ResName = res[i].ResName
			childTree.Res.ResCode = res[i].ResCode
			childTree.Res.Pid = res[i].Pid
			this.toTree(res, childTree, myRes)
			if tree.Children == nil {
				tree.Children = []*models.TreeBean{}
			}
			tree.Children = append(tree.Children, childTree)
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
}
