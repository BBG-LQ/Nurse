package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ServiceOrderModel = (*customServiceOrderModel)(nil)

type (
	// ServiceOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customServiceOrderModel.
	ServiceOrderModel interface {
		serviceOrderModel
	}

	customServiceOrderModel struct {
		*defaultServiceOrderModel
	}
)

// NewServiceOrderModel returns a model for the database table.
func NewServiceOrderModel(conn sqlx.SqlConn, c cache.CacheConf) ServiceOrderModel {
	return &customServiceOrderModel{
		defaultServiceOrderModel: newServiceOrderModel(conn, c),
	}
}
