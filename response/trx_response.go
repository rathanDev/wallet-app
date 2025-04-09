package response

type TrxResponse struct { // Deposit or Withdrawal or Transfer
	TransactionId  string
	WalletId       string
	Amount         uint
	CurrentBalance uint
}
