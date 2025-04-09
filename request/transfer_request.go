package request

type TransferReq struct {
	Amount               uint   `json:"amount" binding:"required"`
	CounterpartyWalletId string `json:"counterpartyWalletId" binding:"required"`
}
