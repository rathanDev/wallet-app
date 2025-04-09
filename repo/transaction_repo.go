package repo

import (
	"wallet-app/entity"

	"gorm.io/gorm"
)

type ITrxRepo interface {
	FindAllTrxs() []entity.TrxEntity
	FindTransactionsByWalletId(walletId string) []entity.TrxEntity
	SaveTrx(trx entity.TrxEntity) error
	SaveTrxWithDbTx(trx entity.TrxEntity, dbTx *gorm.DB) error
	SaveTrxs(trxs []entity.TrxEntity) error
	SaveTrxsWithDbTx(trxs []entity.TrxEntity, dbTx *gorm.DB) error
	DeleteAllTrxs() error
}

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) ITrxRepo {
	return &TransactionRepo{db: db}
}

func (t *TransactionRepo) FindAllTrxs() []entity.TrxEntity {
	var trxs []entity.TrxEntity
	t.db.Find(&trxs)
	return trxs
}

func (t *TransactionRepo) FindTransactionsByWalletId(walletId string) []entity.TrxEntity {
	var transactions []entity.TrxEntity
	t.db.Where("wallet_id = ?", walletId).Find(&transactions)
	return transactions
}

func (t *TransactionRepo) SaveTrx(trx entity.TrxEntity) error {
	return t.db.Save(&trx).Error
}

func (t *TransactionRepo) SaveTrxWithDbTx(trx entity.TrxEntity, tx *gorm.DB) error {
	return tx.Save(&trx).Error
}

func (t *TransactionRepo) SaveTrxs(trxs []entity.TrxEntity) error {
	return t.db.Save(&trxs).Error
}

func (t *TransactionRepo) SaveTrxsWithDbTx(trxs []entity.TrxEntity, tx *gorm.DB) error {
	return tx.Save(&trxs).Error
}

func (t *TransactionRepo) DeleteAllTrxs() error {
	return t.db.Exec("delete from transactions").Error
}
