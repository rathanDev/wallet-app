package db

import (
	"fmt"
	"wallet-app/config"
	"wallet-app/entity"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func InitDb(cfg *config.DatabaseConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.WalletEntity{})
	db.AutoMigrate(&entity.TrxEntity{})

	return db, nil
}
