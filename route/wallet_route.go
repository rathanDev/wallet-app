package route

import (
	"wallet-app/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, controller *controller.WalletController) {
	r.POST("/wallets", controller.CreateWallet)
	r.GET("/wallets/user/:userId", controller.GetWalletsByUserId)

	walletRoute := r.Group("/wallets/:walletId")
	walletRoute.POST("/deposit", controller.DepositMoney)
	walletRoute.POST("/withdraw", controller.WithdrawMoney)
	walletRoute.POST("/transfer", controller.TransferMoney)
	walletRoute.GET("/balance", controller.GetBalance)
	walletRoute.GET("/transactions", controller.GetTransactions)

	r.DELETE("/delete-all", controller.DeleteAll)
}
