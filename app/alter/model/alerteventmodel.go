package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AlertEventModel = (*customAlertEventModel)(nil)

type (
	// AlertEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAlertEventModel.
	AlertEventModel interface {
		alertEventModel
	}

	customAlertEventModel struct {
		*defaultAlertEventModel
	}
)

// NewAlertEventModel returns a model for the database table.
func NewAlertEventModel(conn sqlx.SqlConn, c cache.CacheConf) AlertEventModel {
	return &customAlertEventModel{
		defaultAlertEventModel: newAlertEventModel(conn, c),
	}
}
