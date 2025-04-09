package main

import (
	"fmt"
	"wallet-app/config"
	"wallet-app/controller"
	"wallet-app/db"
	"wallet-app/manager"
	"wallet-app/mapper"
	"wallet-app/repo"
	"wallet-app/route"
	"wallet-app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.Info("Prepare server")

	log.Info("Load config")
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Error("Err loading config; ", err)
		return
	}

	log.Info("Connect to db")
	db, err := db.InitDb(&appConfig.Database)
	if err != nil {
		log.Error("Err connecting to db; ", err)
		return
	}

	dbTxManager := manager.NewDbTxManager(db)
	walletRepo := repo.NewWalletRepo(db)
	transactionRepo := repo.NewTransactionRepo(db)
	mapper := mapper.NewAppMapper()
	service := service.NewWalletService(log, walletRepo, transactionRepo, mapper, dbTxManager)
	controller := controller.NewWalletController(log, service)
	r := gin.Default()
	route.InitRoutes(r, controller)

	serverPort := fmt.Sprintf(":%d", appConfig.Server.Port)
	log.Infof("Start server; port:%s", serverPort)
	r.Run(serverPort)
}
