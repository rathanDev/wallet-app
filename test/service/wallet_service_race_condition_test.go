package service_test

import (
	"testing"
	"time"
	"wallet-app/apperror"
	"wallet-app/entity"
	"wallet-app/manager"
	"wallet-app/mapper"
	"wallet-app/repo"
	"wallet-app/request"
	"wallet-app/response"
	"wallet-app/service"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestWithdrawMoney_RaceCondition(t *testing.T) {
	log := logrus.New()
	log.SetReportCaller(true)

	db, err := gorm.Open(sqlite.Open("file:app_test.db?cache=shared&_journal_mode=WAL&_busy_timeout=5000"), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	require.NoError(t, err)

	db.Migrator().DropTable(&entity.WalletEntity{}, &entity.TrxEntity{})
	require.NoError(t, db.AutoMigrate(&entity.WalletEntity{}, &entity.TrxEntity{}))

	walletId := "wallet123"
	initialBalance := uint(20000)
	withdrawalAmount := uint(10000)
	expectedAmountAfterWithdrawals := uint(0)

	wallet := entity.WalletEntity{
		ID:      walletId,
		Balance: initialBalance,
	}
	require.NoError(t, db.Create(&wallet).Error)

	service := service.NewWalletService(
		log,
		repo.NewWalletRepo(db),
		repo.NewTransactionRepo(db),
		&mapper.AppMapper{},
		manager.NewDbTxManager(db),
	)

	const count = 10
	results := make(chan response.ResonseWrapper, count)

	for i := 0; i < count; i++ {
		result := service.WithdrawMoney(walletId, request.TrxReq{Amount: withdrawalAmount})
		results <- result
	}
	close(results)

	time.Sleep(time.Second * 2)
	walletBalance := service.GetBalance(walletId)
	assert.Equal(t, expectedAmountAfterWithdrawals, walletBalance.Data.(entity.WalletEntity).Balance)

	successCount := 0
	insufficientAmountErrs := 0
	otherErrs := 0
	for res := range results {
		if res.Err.Code == 0 {
			successCount++
		} else if res.Err.Message == apperror.ErrInsufficientAmount.Message {
			insufficientAmountErrs++
		} else {
			otherErrs++
		}
	}

	assert.Equal(t, 2, successCount)
	assert.Equal(t, count-2, insufficientAmountErrs)
	assert.Equal(t, 0, otherErrs)
}
