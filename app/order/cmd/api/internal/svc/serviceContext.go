// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"Nurse/app/order/cmd/api/internal/config"
	"Nurse/app/order/cmd/rpc/orderClient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	OrderRpc orderClient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		OrderRpc: orderClient.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
	}
}
