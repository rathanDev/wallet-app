package service

import (
	"errors"
	"time"

	"wallet-app/apperror"
	"wallet-app/common"
	"wallet-app/entity"
	"wallet-app/manager"
	"wallet-app/mapper"
	"wallet-app/repo"
	"wallet-app/request"
	"wallet-app/response"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IWalletService interface {
	CreateWallet(req request.CreateWalletReq) response.ResonseWrapper
	GetWalletsByUserId(userId string) response.ResonseWrapper

	DepositMoney(walletId string, req request.TrxReq) response.ResonseWrapper
	WithdrawMoney(walletId string, req request.TrxReq) response.ResonseWrapper
	TransferMoney(walletId string, req request.TransferReq) response.ResonseWrapper

	GetBalance(walletId string) response.ResonseWrapper
	GetTransactions(walletId string) response.ResonseWrapper

	DeleteAll() response.ResonseWrapper
	GetAllTrxs() response.ResonseWrapper
}

type WalletService struct {
	log         *logrus.Logger
	dbTxManager manager.IDbTxManager
	// userRepo        repo.IUserRepo
	walletRepo repo.IWalletRepo
	trxRepo    repo.ITrxRepo
	mapper     *mapper.AppMapper
}

func NewWalletService(log *logrus.Logger, walletRepo repo.IWalletRepo, trxRepo repo.ITrxRepo, mapper *mapper.AppMapper, dbTxManager manager.IDbTxManager) IWalletService {
	return &WalletService{log: log, walletRepo: walletRepo, trxRepo: trxRepo, mapper: mapper, dbTxManager: dbTxManager}
}

func (w *WalletService) CreateWallet(req request.CreateWalletReq) response.ResonseWrapper {
	w.log.Infof("CreateWallet; req:%v", req)

	wallet := entity.WalletEntity{ID: uuid.New().String(), UserId: req.UserId, Balance: 0, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if err := w.walletRepo.SaveWallet(wallet); err != nil {
		w.log.Error("Err saving wallet; ", err)
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}

	return response.ResonseWrapper{Data: wallet}
}

func (w *WalletService) GetWalletsByUserId(userId string) response.ResonseWrapper {
	w.log.Infof("GetWalletsByUserId; userId:%s", userId)
	wallets := w.walletRepo.FindWalletsByUserId(userId)
	w.log.Info("Wallets ", wallets)
	return response.ResonseWrapper{Data: wallets}
}

func (w *WalletService) DepositMoney(walletId string, req request.TrxReq) response.ResonseWrapper {
	w.log.Infof("DepositMoney; walletId:%s", walletId)

	dbTx := w.dbTxManager.GetTx().Begin()
	if dbTx.Error != nil {
		w.log.Error("Failed creating dbTrx ", dbTx.Error)
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}
	defer func() {
		if r := recover(); r != nil {
			dbTx.Rollback()
		}
	}()

	wallet, err := w.walletRepo.FindWalletByIdWithTx(walletId, dbTx.Clauses(clause.Locking{Strength: "UPDATE"}))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.log.Errorf("Wallet not found; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrWalletNotFound}
	}
	w.log.Info("Wallet ", wallet)

	wallet.Balance += req.Amount
	if err := w.walletRepo.SaveWalletWithTx(wallet, dbTx); err != nil {
		w.log.Errorf("Err saving wallet; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}

	trx := entity.TrxEntity{ID: uuid.New().String(), WalletId: wallet.ID, Amount: req.Amount, TrxType: common.TrxTypeDeposit, CreatedAt: time.Now()}
	w.log.Info("trx ", trx)
	if err := w.trxRepo.SaveTrxWithDbTx(trx, dbTx); err != nil {
		w.log.Errorf("Err saving trx; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}

	if err := dbTx.Commit().Error; err != nil {
		w.log.Error("Err at commit ", err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}
	w.log.Info("Done committing")

	return response.ResonseWrapper{Data: w.mapper.ToTrxResponse(trx, wallet.Balance)}
}

func (w *WalletService) WithdrawMoney(walletId string, req request.TrxReq) response.ResonseWrapper {
	w.log.Infof("WithdrawMoney; walletId:%s", walletId)

	dbTx := w.dbTxManager.GetTx().Begin()
	if dbTx.Error != nil {
		w.log.Error("Failed creating dbTrx ", dbTx.Error)
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}
	defer func() {
		if r := recover(); r != nil {
			dbTx.Rollback()
		}
	}()
	w.log.Info("DbTrx created")

	wallet, err := w.walletRepo.FindWalletByIdWithTx(walletId, dbTx.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.log.Errorf("Wallet not found; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrWalletNotFound}
	}
	w.log.Info("Wallet ", wallet)

	if wallet.Balance < req.Amount {
		w.log.Errorf("Insufficient amount; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInsufficientAmount}
	}
	wallet.Balance -= req.Amount
	if err := w.walletRepo.SaveWalletWithTx(wallet, dbTx); err != nil {
		w.log.Errorf("Err saving wallet; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}

	trx := entity.TrxEntity{ID: uuid.New().String(), WalletId: wallet.ID, Amount: req.Amount, TrxType: common.TrxTypeWithdrawal, CreatedAt: time.Now()}
	w.log.Info("trx ", trx)
	if err := w.trxRepo.SaveTrxWithDbTx(trx, dbTx); err != nil {
		w.log.Errorf("Err saving trx; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}

	if err := dbTx.Commit().Error; err != nil {
		w.log.Error("Err at commit ", err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}
	w.log.Info("Done committing")

	return response.ResonseWrapper{Data: w.mapper.ToTrxResponse(trx, wallet.Balance)}
}

func (w *WalletService) TransferMoney(walletId string, req request.TransferReq) response.ResonseWrapper {
	w.log.Infof("TransferMoney; walletId:%s", walletId)

	dbTx := w.dbTxManager.GetTx().Begin()
	if dbTx.Error != nil {
		w.log.Error("Failed creating dbTrx ", dbTx.Error)
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}
	defer func() {
		if r := recover(); r != nil {
			dbTx.Rollback()
		}
	}()

	wallet, err := w.walletRepo.FindWalletByIdWithTx(walletId, dbTx.Clauses(clause.Locking{Strength: "UPDATE"}))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.log.Errorf("Wallet not found; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrWalletNotFound}
	}
	w.log.Info("Wallet ", wallet)
	if wallet.Balance < req.Amount {
		w.log.Errorf("Insufficient amount; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInsufficientAmount}
	}

	if walletId == req.CounterpartyWalletId {
		w.log.Errorf("CounterpartyWalletId same as walletId; walletId:%s counterpartyWalletId:%s", walletId, req.CounterpartyWalletId)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrCounterpartyWalletCannotBeSameAsUserWallet}
	}

	counterpartyWallet, err := w.walletRepo.FindWalletByIdWithTx(req.CounterpartyWalletId, dbTx.Clauses(clause.Locking{Strength: "UPDATE"}))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.log.Errorf("CounterpartyWallet not found; walletId:%s %v", walletId, err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrCounterpartyWalletNotFound}
	}
	w.log.Info("CounterpartyWallet ", counterpartyWallet)

	wallet.Balance -= req.Amount
	counterpartyWallet.Balance += req.Amount
	wallets := []entity.WalletEntity{wallet, counterpartyWallet}

	if err := w.walletRepo.SaveWalletsWithTx(wallets, dbTx); err != nil {
		w.log.Error("Err saving wallets; ", err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}

	groupId := uuid.New().String()
	trx := entity.TrxEntity{ID: uuid.New().String(), WalletId: wallet.ID, Amount: req.Amount, TrxType: common.TrxTypeTransferOut, GroupId: groupId, CreatedAt: time.Now()}
	counterpartyTrx := entity.TrxEntity{ID: uuid.New().String(), WalletId: wallet.ID, Amount: req.Amount, TrxType: common.TrxTypeTransferIn, GroupId: groupId, CreatedAt: time.Now()}
	trxs := []entity.TrxEntity{trx, counterpartyTrx}
	if err := w.trxRepo.SaveTrxsWithDbTx(trxs, dbTx); err != nil {
		w.log.Error("Err saving trxs; ", err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}
	w.log.Info("Trxs ", trxs)

	if err := dbTx.Commit().Error; err != nil {
		w.log.Error("Err at commit ", err)
		dbTx.Rollback()
		return response.ResonseWrapper{Err: apperror.ErrInternalServer}
	}
	w.log.Info("Done committing")

	trxRes := w.mapper.ToTrxResponse(trx, wallet.Balance)
	return response.ResonseWrapper{Data: trxRes}
}

func (w *WalletService) GetBalance(walletId string) response.ResonseWrapper {
	w.log.Infof("GetBalance; walletId:%s", walletId)
	wallet, err := w.walletRepo.FindWalletById(walletId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.log.Errorf("Wallet not found; walletId:%s %v", walletId, err)
		return response.ResonseWrapper{Err: apperror.ErrWalletNotFound}
	}
	w.log.Info("Wallet ", wallet)
	return response.ResonseWrapper{Data: wallet}
}

func (w *WalletService) GetTransactions(walletId string) response.ResonseWrapper {
	w.log.Infof("GetTransactions; walletId:%s", walletId)
	wallet, err := w.walletRepo.FindWalletById(walletId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.log.Errorf("Wallet not found; walletId:%s %v", walletId, err)
		return response.ResonseWrapper{Err: apperror.ErrWalletNotFound}
	}
	w.log.Info("Wallet ", wallet)
	trxs := w.trxRepo.FindTransactionsByWalletId(walletId)
	w.log.Info("Trxs ", trxs)
	return response.ResonseWrapper{Data: trxs}
}

func (w *WalletService) GetAllWallets() response.ResonseWrapper {
	w.log.Info("GetAllWallets")
	wallets := w.walletRepo.FindAllWallets()
	w.log.Info("Wallets ", wallets)
	return response.ResonseWrapper{Data: wallets}
}

func (w *WalletService) GetAllTrxs() response.ResonseWrapper {
	w.log.Info("GetAllTrxs")
	trxs := w.trxRepo.FindAllTrxs()
	w.log.Info("Trxs ", trxs)
	return response.ResonseWrapper{Data: trxs}
}

func (w *WalletService) DeleteAll() response.ResonseWrapper {
	w.log.Info("DeleteAll")
	w.walletRepo.DeleteAllWallets()
	w.trxRepo.DeleteAllTrxs()
	return response.ResonseWrapper{}
}
