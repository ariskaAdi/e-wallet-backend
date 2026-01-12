package wallet

import "time"

type WalletResponse struct {
	UserPublicId string    `json:"user_public_id"`
	Balance      int64     `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}