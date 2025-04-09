package response

import "wallet-app/apperror"

type ResonseWrapper struct {
	Data interface{}
	Err  apperror.AppError
}