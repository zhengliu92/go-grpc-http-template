package myservicelogic

import (
	"context"

	"go-grpc-http-template/internal/svc"
	"go-grpc-http-template/internal/util"
	"go-grpc-http-template/myService"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *myService.LoginRequest) (*myService.LoginResponse, error) {
	// TODO: 替换为真实的数据库查询和密码验证
	if in.Username != "admin" || in.Password != "admin123" {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, expire, err := util.GenerateToken(
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
		1,           // userId（从数据库获取）
		in.Username,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &myService.LoginResponse{
		Token:  token,
		Expire: expire,
	}, nil
}
