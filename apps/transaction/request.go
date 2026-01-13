package transaction

type CreateTransactionRequestPayload struct {
	WalletId     int64  `json:"wallet_id"`
	UserPublicId string `json:"user_public_id"`
	Amount       int64  `json:"amount"`
}