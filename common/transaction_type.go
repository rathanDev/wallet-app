package common

type TrxType string

const (
	TrxTypeDeposit     TrxType = "deposit"
	TrxTypeWithdrawal  TrxType = "withdrawal"
	TrxTypeTransferIn  TrxType = "transfer_in"
	TrxTypeTransferOut TrxType = "transfer_out"
)
