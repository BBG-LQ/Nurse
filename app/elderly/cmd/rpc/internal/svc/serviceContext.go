package svc

import (
	"Nurse/app/elderly/cmd/rpc/internal/config"
	"Nurse/app/elderly/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	RedisClient *redis.Redis

	ElderlyArchiveModel    model.ElderlyArchiveModel
	ElderlyMedicationModel model.ElderlyMedicationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config: c,

		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),

		ElderlyArchiveModel:    model.NewElderlyArchiveModel(SqlConn, c.Cache),
		ElderlyMedicationModel: model.NewElderlyMedicationModel(SqlConn, c.Cache),
	}
}
