package controller

import (
	"net/http"

	"wallet-app/apperror"
	"wallet-app/request"
	"wallet-app/response"
	"wallet-app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	log     *logrus.Logger
	service service.IWalletService
}

func NewWalletController(log *logrus.Logger, service service.IWalletService) *WalletController {
	return &WalletController{log: log, service: service}
}

func (w *WalletController) CreateWallet(c *gin.Context) {
	var req request.CreateWalletReq
	if err := c.ShouldBindJSON(&req); err != nil {
		w.log.Error("Err ", err.Error())
		appError := apperror.AppError{Code: 400, Message: err.Error()}
		c.JSON(http.StatusBadRequest, response.ResonseWrapper{Err: appError})
		return
	}
	res := w.service.CreateWallet(req)
	if res.Err.Code != 0 {
		c.JSON(res.Err.Code, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (w *WalletController) GetWalletsByUserId(c *gin.Context) {
	res := w.service.GetWalletsByUserId(c.Param("userId"))
	c.JSON(http.StatusOK, res)
}

func (w *WalletController) DepositMoney(c *gin.Context) {
	var req request.TrxReq
	if err := c.ShouldBindJSON(&req); err != nil {
		w.log.Error("Err ", err.Error())
		appError := apperror.AppError{Code: 400, Message: err.Error()}
		c.JSON(http.StatusBadRequest, response.ResonseWrapper{Err: appError})
		return
	}
	res := w.service.DepositMoney(c.Param("walletId"), req)
	if res.Err.Code != 0 {
		c.JSON(res.Err.Code, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (w *WalletController) WithdrawMoney(c *gin.Context) {
	var req request.TrxReq
	if err := c.ShouldBindJSON(&req); err != nil {
		w.log.Error("Err ", err.Error())
		appError := apperror.AppError{Code: 400, Message: err.Error()}
		c.JSON(http.StatusBadRequest, response.ResonseWrapper{Err: appError})
		return
	}
	res := w.service.WithdrawMoney(c.Param("walletId"), req)
	if res.Err.Code != 0 {
		c.JSON(res.Err.Code, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (w *WalletController) TransferMoney(c *gin.Context) {
	var req request.TransferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		w.log.Error("Err ", err.Error())
		appError := apperror.AppError{Code: 400, Message: err.Error()}
		c.JSON(http.StatusBadRequest, response.ResonseWrapper{Err: appError})
		return
	}
	res := w.service.TransferMoney(c.Param("walletId"), req)
	if res.Err.Code != 0 {
		c.JSON(res.Err.Code, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (w *WalletController) GetBalance(c *gin.Context) {
	res := w.service.GetBalance(c.Param("walletId"))
	if res.Err.Code != 0 {
		c.JSON(res.Err.Code, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (w *WalletController) GetTransactions(c *gin.Context) {
	res := w.service.GetTransactions(c.Param("walletId"))
	if res.Err.Code != 0 {
		c.JSON(res.Err.Code, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (w *WalletController) DeleteAll(c *gin.Context) {
	res := w.service.DeleteAll()
	c.JSON(http.StatusOK, res)
}
