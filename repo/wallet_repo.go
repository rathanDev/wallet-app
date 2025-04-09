package repo

import (
	"wallet-app/entity"

	"gorm.io/gorm"
)

type IWalletRepo interface {
	FindWalletById(walletId string) (entity.WalletEntity, error)
	FindWalletsByUserId(userId string) []entity.WalletEntity
	FindWalletByIdWithTx(walletId string, tx *gorm.DB) (entity.WalletEntity, error)
	FindAllWallets() []entity.WalletEntity
	SaveWallet(wallet entity.WalletEntity) error
	SaveWalletWithTx(wallet entity.WalletEntity, tx *gorm.DB) error
	SaveWallets(wallets []entity.WalletEntity) error
	SaveWalletsWithTx(wallets []entity.WalletEntity, tx *gorm.DB) error
	DeleteAllWallets() error
}

type WalletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) IWalletRepo {
	return &WalletRepo{db: db}
}

func (w *WalletRepo) FindWalletById(id string) (entity.WalletEntity, error) {
	var wallet entity.WalletEntity
	err := w.db.Where("id = ?", id).First(&wallet).Error
	return wallet, err
}

func (w *WalletRepo) FindWalletsByUserId(userId string) []entity.WalletEntity {
	var wallets []entity.WalletEntity
	w.db.Where("user_id = ?", userId).Find(&wallets)
	return wallets
}

func (w *WalletRepo) FindWalletByIdWithTx(id string, tx *gorm.DB) (entity.WalletEntity, error) {
	var wallet entity.WalletEntity
	err := tx.Where("id = ?", id).First(&wallet).Error
	return wallet, err
}

func (w *WalletRepo) FindAllWallets() []entity.WalletEntity {
	var wallets []entity.WalletEntity
	w.db.Find(&wallets)
	return wallets
}

func (w *WalletRepo) SaveWallet(wallet entity.WalletEntity) error {
	return w.db.Save(&wallet).Error
}

func (w *WalletRepo) SaveWalletWithTx(wallet entity.WalletEntity, tx *gorm.DB) error {
	return tx.Save(&wallet).Error
}

func (w *WalletRepo) SaveWallets(wallets []entity.WalletEntity) error {
	return w.db.Save(&wallets).Error
}

func (w *WalletRepo) SaveWalletsWithTx(wallets []entity.WalletEntity, tx *gorm.DB) error {
	return tx.Save(&wallets).Error
}

func (w *WalletRepo) DeleteAllWallets() error {
	return w.db.Exec("delete from wallets").Error
}
