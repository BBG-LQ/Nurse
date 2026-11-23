package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VitalRecordModel = (*customVitalRecordModel)(nil)

type (
	// VitalRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVitalRecordModel.
	VitalRecordModel interface {
		vitalRecordModel
	}

	customVitalRecordModel struct {
		*defaultVitalRecordModel
	}
)

// NewVitalRecordModel returns a model for the database table.
func NewVitalRecordModel(conn sqlx.SqlConn, c cache.CacheConf) VitalRecordModel {
	return &customVitalRecordModel{
		defaultVitalRecordModel: newVitalRecordModel(conn, c),
	}
}
