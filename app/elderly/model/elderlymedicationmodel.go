package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ElderlyMedicationModel = (*customElderlyMedicationModel)(nil)

type (
	// ElderlyMedicationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customElderlyMedicationModel.
	ElderlyMedicationModel interface {
		elderlyMedicationModel
	}

	customElderlyMedicationModel struct {
		*defaultElderlyMedicationModel
	}
)

// NewElderlyMedicationModel returns a model for the database table.
func NewElderlyMedicationModel(conn sqlx.SqlConn, c cache.CacheConf) ElderlyMedicationModel {
	return &customElderlyMedicationModel{
		defaultElderlyMedicationModel: newElderlyMedicationModel(conn, c),
	}
}
