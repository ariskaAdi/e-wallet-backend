package topup

import "time"

type WalletEntity struct {
	Id             int       `db:"id"`
	WalletPublicId string    `db:"wallet_public_id"`
	UserPublicId   string    `db:"user_public_id"`
	Name           string    `db:"name"`
	Balance        int64     `db:"balance"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}