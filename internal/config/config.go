package config

import (
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

type JwtAuth struct {
	AccessSecret string
	AccessExpire int64 // token 有效期（秒）
}

type Config struct {
	zrpc.RpcServerConf
	Gateway gateway.GatewayConf
	JwtAuth JwtAuth
}
