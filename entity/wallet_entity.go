package entity

import "time"

type WalletEntity struct {
	ID        string    `gorm:"primaryKey;column:id"`
	UserId    string    `gorm:"column:user_id"`
	Balance   uint      `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (WalletEntity) TableName() string {
	return "wallets"
}
