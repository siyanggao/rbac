package main

import (
	"fmt"
	"log"
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
	orm.Debug = true
	models.RegisterDB()
	beego.Run()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8081))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	beego.Informational("listen 8081 ok")
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterDepartServer(grpcServer, new(rpc.DepartRpc))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("failed to server: %v", err)
	}
	log.Fatal("server 8081 ok")
}
