package myservicelogic

import (
	"context"
	"strconv"

	"go-grpc-http-template/internal/svc"
	"go-grpc-http-template/myService"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *myService.UserInfoRequest) (*myService.UserInfoResponse, error) {
	// JWT 中间件已验证 token，并通过 Grpc-Metadata-* 头将 claims 注入到 gRPC metadata
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "missing metadata")
	}

	userIdVals := md.Get("user-id")
	if len(userIdVals) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing user identity")
	}

	userId, _ := strconv.ParseInt(userIdVals[0], 10, 64)

	username := ""
	if usernameVals := md.Get("username"); len(usernameVals) > 0 {
		username = usernameVals[0]
	}

	return &myService.UserInfoResponse{
		UserId:   userId,
		Username: username,
	}, nil
}
