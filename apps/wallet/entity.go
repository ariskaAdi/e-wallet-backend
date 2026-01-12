package wallet

import (
	"time"
)

type WalletEntity struct {
	Id           int       `db:"id"`
	UserPublicId string `db:"user_public_id"`
	Balance      int64     `db:"balance"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

