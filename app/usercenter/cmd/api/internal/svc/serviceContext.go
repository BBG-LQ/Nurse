// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Nurse/app/usercenter/cmd/api/internal/config"
	"Nurse/app/usercenter/cmd/rpc/usercenterClient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UsercenterRpc usercenterClient.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UsercenterRpc: usercenterClient.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}
