package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SatisfactionModel = (*customSatisfactionModel)(nil)

type (
	// SatisfactionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSatisfactionModel.
	SatisfactionModel interface {
		satisfactionModel
	}

	customSatisfactionModel struct {
		*defaultSatisfactionModel
	}
)

// NewSatisfactionModel returns a model for the database table.
func NewSatisfactionModel(conn sqlx.SqlConn, c cache.CacheConf) SatisfactionModel {
	return &customSatisfactionModel{
		defaultSatisfactionModel: newSatisfactionModel(conn, c),
	}
}
