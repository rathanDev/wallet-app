package mapper

import (
	"wallet-app/entity"
	"wallet-app/response"
)

type AppMapper struct {
}

func NewAppMapper() *AppMapper {
	return &AppMapper{}
}

func (a *AppMapper) ToTrxResponse(e entity.TrxEntity, balance uint) response.TrxResponse {
	return response.TrxResponse{TransactionId: e.ID, WalletId: e.WalletId, Amount: e.Amount, CurrentBalance: balance}
}
