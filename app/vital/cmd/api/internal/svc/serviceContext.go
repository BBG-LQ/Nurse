// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Nurse/app/vital/cmd/api/internal/config"
	"Nurse/app/vital/cmd/rpc/vitalClient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	VitalRpc vitalClient.Vital
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		VitalRpc: vitalClient.NewVital(zrpc.MustNewClient(c.VitalRpcConf)),
	}
}
