package main

import (
	"fmt"
	"net"
	"rbac/models"
	pb "rbac/proto"
	_ "rbac/routers"
	"rbac/rpc"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8081))
	if err != nil {
		beego.Error("failed to listen: %v", err)
	}
	beego.Informational("listen 8081 ok")
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterDepartServer(grpcServer, new(rpc.DepartRpc))
	go grpcServer.Serve(lis)

	beego.Informational("server 8081 ok")

	orm.Debug = true
	models.RegisterDB()
	beego.Run()
	beego.Informational("test2")

}
