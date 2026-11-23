package svc

import (
	"Nurse/app/vital/cmd/rpc/internal/config"
	"Nurse/app/vital/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	VitalRecordModel model.VitalRecordModel

	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	SqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		VitalRecordModel: model.NewVitalRecordModel(SqlConn, c.Cache),
	}
}
