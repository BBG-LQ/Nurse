package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdverseEventModel = (*customAdverseEventModel)(nil)

type (
	// AdverseEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdverseEventModel.
	AdverseEventModel interface {
		adverseEventModel
	}

	customAdverseEventModel struct {
		*defaultAdverseEventModel
	}
)

// NewAdverseEventModel returns a model for the database table.
func NewAdverseEventModel(conn sqlx.SqlConn, c cache.CacheConf) AdverseEventModel {
	return &customAdverseEventModel{
		defaultAdverseEventModel: newAdverseEventModel(conn, c),
	}
}
