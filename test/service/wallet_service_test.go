package service_test

import (
	"testing"
	"wallet-app/apperror"
	"wallet-app/entity"
	"wallet-app/mapper"
	"wallet-app/request"
	"wallet-app/response"
	"wallet-app/service"
	mock_test "wallet-app/test/mock"

	"github.com/glebarez/sqlite"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateWallet(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	const userId = "jana"
	req := request.CreateWalletReq{UserId: userId}

	mockWalletRepo.On("SaveWallet", mock.Anything, mock.Anything).Return(nil)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.CreateWallet(req)
	assert.Equal(t, 0, result.Err.Code)
	assert.Equal(t, userId, result.Data.(entity.WalletEntity).UserId)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestGetWalletsByUserId_found(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	const userIdJana = "jana"

	wallet1 := entity.WalletEntity{ID: "walletId1", UserId: userIdJana}
	wallet2 := entity.WalletEntity{ID: "walletId2", UserId: userIdJana}
	wallets := []entity.WalletEntity{wallet1, wallet2}

	mockWalletRepo.On("FindWalletsByUserId", userIdJana).Return(wallets)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.GetWalletsByUserId(userIdJana)
	assert.Equal(t, 0, result.Err.Code)
	assert.Equal(t, 2, len(result.Data.([]entity.WalletEntity)))
	assert.Equal(t, userIdJana, result.Data.([]entity.WalletEntity)[0].UserId)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestGetWalletsByUserId_notFound(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	const userIdNone = "none"

	wallets := []entity.WalletEntity{}

	mockWalletRepo.On("FindWalletsByUserId", userIdNone).Return(wallets)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.GetWalletsByUserId(userIdNone)
	assert.Equal(t, 0, result.Err.Code)
	assert.Equal(t, 0, len(result.Data.([]entity.WalletEntity)))
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestDepositMoney_walletNotFound(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	walletId := "walletIdNot"
	amount := uint(1000)
	req := request.TrxReq{Amount: amount}

	mockTxManager.On("GetTx").Return(getTestDB(t))
	mockWalletRepo.On("FindWalletByIdWithTx", walletId, mock.Anything).Return(entity.WalletEntity{}, gorm.ErrRecordNotFound)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.DepositMoney(walletId, req)

	assert.Equal(t, 400, result.Err.Code)
	assert.Equal(t, apperror.ErrWalletNotFound.Message, result.Err.Message)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestDepositMoney_Success(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	walletId := "wallet123"
	amount := uint(1000)
	req := request.TrxReq{Amount: amount}
	wallet := entity.WalletEntity{ID: walletId, Balance: 5000}

	mockTxManager.On("GetTx").Return(getTestDB(t))
	mockWalletRepo.On("FindWalletByIdWithTx", walletId, mock.Anything).Return(wallet, nil)
	mockWalletRepo.On("SaveWalletWithTx", mock.Anything, mock.Anything).Return(nil)
	mockTrxRepo.On("SaveTrxWithDbTx", mock.Anything, mock.Anything).Return(nil)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.DepositMoney(walletId, req)

	assert.Equal(t, 0, result.Err.Code)
	assert.Equal(t, walletId, result.Data.(response.TrxResponse).WalletId)
	assert.Equal(t, uint(6000), result.Data.(response.TrxResponse).CurrentBalance)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestWithdrawMoney_InsufficientAmount(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	walletId := "wallet123"
	req := request.TrxReq{Amount: uint(10000)}
	wallet := entity.WalletEntity{ID: walletId, Balance: 5000}

	mockTxManager.On("GetTx").Return(getTestDB(t))
	mockWalletRepo.On("FindWalletByIdWithTx", walletId, mock.Anything).Return(wallet, nil)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.WithdrawMoney(walletId, req)

	assert.Equal(t, 400, result.Err.Code)
	assert.Equal(t, apperror.ErrInsufficientAmount.Message, result.Err.Message)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestWithdrawMoney_success(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	walletId := "wallet123"
	amount := uint(10000)
	req := request.TrxReq{Amount: amount}
	wallet := entity.WalletEntity{ID: walletId, Balance: 20000}

	mockTxManager.On("GetTx").Return(getTestDB(t))
	mockWalletRepo.On("FindWalletByIdWithTx", walletId, mock.Anything).Return(wallet, nil)
	mockWalletRepo.On("SaveWalletWithTx", mock.Anything, mock.Anything).Return(nil)
	mockTrxRepo.On("SaveTrxWithDbTx", mock.Anything, mock.Anything).Return(nil)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.WithdrawMoney(walletId, req)
	assert.Equal(t, 0, result.Err.Code)
	assert.Equal(t, uint(10000), result.Data.(response.TrxResponse).CurrentBalance)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestTransferMoney_CounterpartyWalletSameAsUserWallet(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	walletId := "wallet_mine"
	counterpartyWalletId := "wallet_mine"
	amount := uint(5000)
	req := request.TransferReq{Amount: amount, CounterpartyWalletId: counterpartyWalletId}
	wallet := entity.WalletEntity{ID: walletId, Balance: 20000}

	mockTxManager.On("GetTx").Return(getTestDB(t))
	mockWalletRepo.On("FindWalletByIdWithTx", walletId, mock.Anything).Return(wallet, nil)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.TransferMoney(walletId, req)

	assert.Equal(t, 400, result.Err.Code)
	assert.Equal(t, apperror.ErrCounterpartyWalletCannotBeSameAsUserWallet.Message, result.Err.Message)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestTransferMoney_CounterpartyWalletNotFound(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	walletId := "wallet_mine"
	counterpartyWalletId := "wallet_counterparty"
	amount := uint(5000)
	req := request.TransferReq{Amount: amount, CounterpartyWalletId: counterpartyWalletId}
	wallet := entity.WalletEntity{ID: walletId, Balance: 20000}

	mockTxManager.On("GetTx").Return(getTestDB(t))
	mockWalletRepo.On("FindWalletByIdWithTx", walletId, mock.Anything).Return(wallet, nil)
	mockWalletRepo.On("FindWalletByIdWithTx", counterpartyWalletId, mock.Anything).Return(entity.WalletEntity{}, gorm.ErrRecordNotFound)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.TransferMoney(walletId, req)

	assert.Equal(t, 400, result.Err.Code)
	assert.Equal(t, apperror.ErrCounterpartyWalletNotFound.Message, result.Err.Message)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func TestTransferMoney_success(t *testing.T) {
	mockWalletRepo := new(mock_test.MockWalletRepo)
	mockTrxRepo := new(mock_test.MockTrxRepo)
	mockTxManager := new(mock_test.MockDbTxManager)

	walletId := "wallet_mine"
	counterpartyWalletId := "wallet_counterparty"
	amount := uint(5000)
	req := request.TransferReq{Amount: amount, CounterpartyWalletId: counterpartyWalletId}
	wallet := entity.WalletEntity{ID: walletId, Balance: 20000}
	counterpartyWallet := entity.WalletEntity{ID: counterpartyWalletId, Balance: 1000}

	mockTxManager.On("GetTx").Return(getTestDB(t))
	mockWalletRepo.On("FindWalletByIdWithTx", walletId, mock.Anything).Return(wallet, nil)
	mockWalletRepo.On("FindWalletByIdWithTx", counterpartyWalletId, mock.Anything).Return(counterpartyWallet, nil)
	mockWalletRepo.On("SaveWalletsWithTx", mock.Anything, mock.Anything).Return(nil)
	mockTrxRepo.On("SaveTrxsWithDbTx", mock.Anything, mock.Anything).Return(nil)

	service := service.NewWalletService(
		logrus.New(),
		mockWalletRepo,
		mockTrxRepo,
		&mapper.AppMapper{},
		mockTxManager,
	)

	result := service.TransferMoney(walletId, req)

	assert.Equal(t, 0, result.Err.Code)
	assert.Equal(t, uint(15000), result.Data.(response.TrxResponse).CurrentBalance)
	mockWalletRepo.AssertExpectations(t)
	mockTrxRepo.AssertExpectations(t)
}

func getTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory DB: %v", err)
	}
	err = db.AutoMigrate(&entity.WalletEntity{}, &entity.TrxEntity{})
	if err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}
	return db
}
