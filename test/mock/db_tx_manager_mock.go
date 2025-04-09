package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDbTxManager struct {
	mock.Mock
}

func (m *MockDbTxManager) GetTx() *gorm.DB {
	args := m.Called()
	if tx := args.Get(0); tx != nil {
		return tx.(*gorm.DB)
	}
	return nil
}
