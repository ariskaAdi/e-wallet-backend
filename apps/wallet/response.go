package wallet

import "time"

type MyWalletResponse struct {
	UserPublicId string    `json:"user_public_id"`
	WalletPublicId string `json:"wallet_public_id"`
	Name         string    `json:"name"`
	Balance      int64     `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FindWalletResponse struct {
	Name string `json:"name"`
}

