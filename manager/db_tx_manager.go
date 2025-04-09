package manager

import "gorm.io/gorm"

type IDbTxManager interface {
	GetTx() *gorm.DB
}

type DbTxManager struct {
	db *gorm.DB
}

func NewDbTxManager(db *gorm.DB) IDbTxManager {
	return &DbTxManager{db: db}
}

func (t *DbTxManager) GetTx() *gorm.DB {
	return t.db
}
