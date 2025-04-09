package request

type TrxReq struct { // Deposit or Withdrawal
	Amount uint `json:"amount" binding:"required"`
}
