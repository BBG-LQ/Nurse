package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FamilyBindingModel = (*customFamilyBindingModel)(nil)

type (
	// FamilyBindingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFamilyBindingModel.
	FamilyBindingModel interface {
		familyBindingModel
	}

	customFamilyBindingModel struct {
		*defaultFamilyBindingModel
	}
)

// NewFamilyBindingModel returns a model for the database table.
func NewFamilyBindingModel(conn sqlx.SqlConn, c cache.CacheConf) FamilyBindingModel {
	return &customFamilyBindingModel{
		defaultFamilyBindingModel: newFamilyBindingModel(conn, c),
	}
}
