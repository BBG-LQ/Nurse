package svc

import (
	"Nurse/app/elderly/cmd/rpc/elderlyClient"
	"Nurse/app/nurse/cmd/rpc/nurseClient"
	"Nurse/app/order/cmd/rpc/internal/config"
	"Nurse/app/order/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis

	SatisfactionModel model.SatisfactionModel
	ServiceOrderModel model.ServiceOrderModel

	NurseRpc   nurseClient.Nurse
	ElderlyRpc elderlyClient.Elderly
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config: c,

		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),

		SatisfactionModel: model.NewSatisfactionModel(SqlConn, c.Cache),
		ServiceOrderModel: model.NewServiceOrderModel(SqlConn, c.Cache),

		NurseRpc:   nurseClient.NewNurse(zrpc.MustNewClient(c.NurseRpcConf)),
		ElderlyRpc: elderlyClient.NewElderly(zrpc.MustNewClient(c.ElderlyRpcConf)),
	}
}
