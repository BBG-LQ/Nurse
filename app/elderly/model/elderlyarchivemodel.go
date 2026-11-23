package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ElderlyArchiveModel = (*customElderlyArchiveModel)(nil)

type (
	// ElderlyArchiveModel is an interface to be customized, add more methods here,
	// and implement the added methods in customElderlyArchiveModel.
	ElderlyArchiveModel interface {
		elderlyArchiveModel
	}

	customElderlyArchiveModel struct {
		*defaultElderlyArchiveModel
	}
)

// NewElderlyArchiveModel returns a model for the database table.
func NewElderlyArchiveModel(conn sqlx.SqlConn, c cache.CacheConf) ElderlyArchiveModel {
	return &customElderlyArchiveModel{
		defaultElderlyArchiveModel: newElderlyArchiveModel(conn, c),
	}
}
