package svc

import (
	"Nurse/app/alter/cmd/rpc/internal/config"
	"Nurse/app/alter/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	RedisClient *redis.Redis

	AdverseEventModel model.AdverseEventModel
	AlertEventModel   model.AlertEventModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		AlertEventModel:   model.NewAlertEventModel(SqlConn, c.Cache),
		AdverseEventModel: model.NewAdverseEventModel(SqlConn, c.Cache),
	}
}
