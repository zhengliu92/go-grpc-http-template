package main

import (
	"flag"
	"fmt"

	"go-grpc-http-template/internal/config"
	"go-grpc-http-template/internal/middleware"
	myserviceServer "go-grpc-http-template/internal/server/myservice"
	"go-grpc-http-template/internal/svc"
	"go-grpc-http-template/myService"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/myService.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// gRPC server
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		myService.RegisterMyServiceServer(grpcServer, myserviceServer.NewMyServiceServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// HTTP gateway (gRPC -> HTTP)，使用响应包装中间件
	gw := gateway.MustNewServer(c.Gateway, gateway.WithMiddleware(middleware.WrapResponse))

	// 使用 ServiceGroup 同时管理两个服务
	// ServiceGroup 内部已处理信号监听 (SIGINT/SIGTERM) 和优雅关闭
	group := service.NewServiceGroup()
	defer group.Stop()
	group.Add(s)
	group.Add(gw)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting http gateway at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	group.Start()
}
