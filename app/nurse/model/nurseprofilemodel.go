package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NurseProfileModel = (*customNurseProfileModel)(nil)

type (
	// NurseProfileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNurseProfileModel.
	NurseProfileModel interface {
		nurseProfileModel
	}

	customNurseProfileModel struct {
		*defaultNurseProfileModel
	}
)

// NewNurseProfileModel returns a model for the database table.
func NewNurseProfileModel(conn sqlx.SqlConn, c cache.CacheConf) NurseProfileModel {
	return &customNurseProfileModel{
		defaultNurseProfileModel: newNurseProfileModel(conn, c),
	}
}
