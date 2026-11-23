// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Nurse/app/nurse/cmd/api/internal/config"
	"Nurse/app/nurse/cmd/rpc/nurseClient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	NurseRpc nurseClient.Nurse
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		NurseRpc: nurseClient.NewNurse(zrpc.MustNewClient(c.NurseRpcConf)),
	}
}
