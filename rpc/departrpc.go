package rpc

import (
	pb "rbac/proto"
	"rbac/services"
	"rbac/utils"

	"github.com/astaxie/beego"
	"golang.org/x/net/context"
)

type DepartRpc struct {
	service services.DepartService
}

func (this *DepartRpc) GetAllChildDepartByUserId(ctx context.Context, in *pb.User) (*pb.ChildDepart, error) {
	child, err := this.service.GetAllChildDepartByUserId(int(in.UserId), in.WithMe)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	childDepart := new(pb.ChildDepart)
	childDepart.Depart = utils.IntToInt32Slice(*child)
	return childDepart, nil
}
