// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Nurse/app/alter/cmd/api/internal/config"
	"Nurse/app/alter/cmd/rpc/alertClient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	AlertRpc alertClient.Alert
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		AlertRpc: alertClient.NewAlert(zrpc.MustNewClient(c.AlertRpcConf)),
	}
}
