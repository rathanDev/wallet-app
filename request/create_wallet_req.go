package request

type CreateWalletReq struct {
	UserId string `json:"userId" binding:"required"`
}
