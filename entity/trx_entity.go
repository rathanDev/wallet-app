package entity

import (
	"time"
	"wallet-app/common"
)

type TrxEntity struct {
	ID                   string         `gorm:"primaryKey;column:id"`
	WalletId             string         `gorm:"column:wallet_id"`
	Amount               uint           `gorm:"column:amount"`
	CounterpartyWalletId string         `gorm:"column:counterparty_wallet_id"`
	TrxType              common.TrxType `gorm:"column:trx_type"`
	GroupId              string         `gorm:"column:group_id"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
}

func (TrxEntity) TableName() string {
	return "transactions"
}
