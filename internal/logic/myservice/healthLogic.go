package myservicelogic

import (
	"context"

	"go-grpc-http-template/internal/svc"
	"go-grpc-http-template/myService"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HealthLogic) Health(in *myService.Request) (*myService.Response, error) {
	return &myService.Response{
		Status: "ok",
	}, nil
}
