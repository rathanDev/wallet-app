package mock_test

import (
	"wallet-app/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockTrxRepo struct {
	mock.Mock
}

func NewMockTrxRepo() *MockTrxRepo {
	return &MockTrxRepo{}
}

func (m *MockTrxRepo) FindAllTrxs() []entity.TrxEntity {
	args := m.Called()
	return args.Get(0).([]entity.TrxEntity)
}

func (m *MockTrxRepo) FindTransactionsByWalletId(walletId string) []entity.TrxEntity {
	args := m.Called(walletId)
	return args.Get(0).([]entity.TrxEntity)
}

func (m *MockTrxRepo) SaveTrx(trx entity.TrxEntity) error {
	args := m.Called(trx)
	return args.Error(0)
}

func (m *MockTrxRepo) SaveTrxWithDbTx(trx entity.TrxEntity, tx *gorm.DB) error {
	args := m.Called(trx, tx)
	return args.Error(0)
}

func (m *MockTrxRepo) SaveTrxs(trxs []entity.TrxEntity) error {
	args := m.Called(trxs)
	return args.Error(0)
}

func (m *MockTrxRepo) SaveTrxsWithDbTx(trxs []entity.TrxEntity, tx *gorm.DB) error {
	args := m.Called(trxs, tx)
	return args.Error(0)
}

func (m *MockTrxRepo) DeleteAllTrxs() error {
	args := m.Called()
	return args.Error(0)
}
