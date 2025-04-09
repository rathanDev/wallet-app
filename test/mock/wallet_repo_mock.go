package mock_test

import (
	"wallet-app/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockWalletRepo struct {
	mock.Mock
}

func NewMockWalletRepo() *MockWalletRepo {
	return &MockWalletRepo{}
}

func (w *MockWalletRepo) FindWalletById(id string) (entity.WalletEntity, error) {
	args := w.Called(id)
	return args.Get(0).(entity.WalletEntity), args.Error(1)
}

func (w *MockWalletRepo) FindWalletsByUserId(userId string) []entity.WalletEntity {
	args := w.Called(userId)
	return args.Get(0).([]entity.WalletEntity)
}

func (w *MockWalletRepo) FindWalletByIdWithTx(id string, tx *gorm.DB) (entity.WalletEntity, error) {
	args := w.Called(id, tx)
	return args.Get(0).(entity.WalletEntity), args.Error(1)
}

func (w *MockWalletRepo) FindAllWallets() []entity.WalletEntity {
	args := w.Called()
	return args.Get(0).([]entity.WalletEntity)
}

func (w *MockWalletRepo) SaveWallet(wallet entity.WalletEntity) error {
	args := w.Called()
	return args.Error(0)
}

func (w *MockWalletRepo) SaveWalletWithTx(wallet entity.WalletEntity, tx *gorm.DB) error {
	args := w.Called(wallet, tx)
	return args.Error(0)
}

func (w *MockWalletRepo) SaveWallets(wallets []entity.WalletEntity) error {
	args := w.Called(wallets)
	return args.Error(0)
}

func (w *MockWalletRepo) SaveWalletsWithTx(wallets []entity.WalletEntity, tx *gorm.DB) error {
	args := w.Called(wallets, tx)
	return args.Error(0)
}

func (w *MockWalletRepo) DeleteAllWallets() error {
	args := w.Called()
	return args.Error(0)
}
