package apperror

type AppError struct {
	Code    int
	Message string
}

var (
	ErrWalletIdNotFound    = AppError{Code: 400, Message: "wallet id not found"}
	ErrIncompatibleRequest = AppError{Code: 400, Message: "incompatible request"}

	ErrUserNotFound                               = AppError{Code: 400, Message: "user not found"}
	ErrWalletNotFound                             = AppError{Code: 400, Message: "wallet not found"}
	ErrCounterpartyWalletNotFound                 = AppError{Code: 400, Message: "counterparty wallet not found"}
	ErrCounterpartyWalletCannotBeSameAsUserWallet = AppError{Code: 400, Message: "counterparty wallet can not be same as user wallet"}
	ErrInsufficientAmount                         = AppError{Code: 400, Message: "insufficient amount"}

	ErrInternalServer = AppError{Code: 500, Message: "internal server error"}
)
