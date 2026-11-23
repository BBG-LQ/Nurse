package svc

import (
	"Nurse/app/nurse/cmd/rpc/internal/config"
	"Nurse/app/nurse/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis

	NurseProfileModel   model.NurseProfileModel
	TrainingRecordModel model.TrainingRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config: c,

		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),

		NurseProfileModel:   model.NewNurseProfileModel(SqlConn, c.Cache),
		TrainingRecordModel: model.NewTrainingRecordModel(SqlConn, c.Cache),
	}
}
