package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TrainingRecordModel = (*customTrainingRecordModel)(nil)

type (
	// TrainingRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTrainingRecordModel.
	TrainingRecordModel interface {
		trainingRecordModel
	}

	customTrainingRecordModel struct {
		*defaultTrainingRecordModel
	}
)

// NewTrainingRecordModel returns a model for the database table.
func NewTrainingRecordModel(conn sqlx.SqlConn, c cache.CacheConf) TrainingRecordModel {
	return &customTrainingRecordModel{
		defaultTrainingRecordModel: newTrainingRecordModel(conn, c),
	}
}
