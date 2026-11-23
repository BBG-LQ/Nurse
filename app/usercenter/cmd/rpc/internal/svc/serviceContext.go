package svc

import (
	"Nurse/app/usercenter/cmd/rpc/internal/config"
	"Nurse/app/usercenter/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	SysUserModel       model.SysUserModel
	FamilyBindingModel model.FamilyBindingModel

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

		SysUserModel:       model.NewSysUserModel(SqlConn, c.Cache),
		FamilyBindingModel: model.NewFamilyBindingModel(SqlConn, c.Cache),
	}
}
