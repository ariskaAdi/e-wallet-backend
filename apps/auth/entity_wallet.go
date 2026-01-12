package auth

import (
	"time"

	"github.com/google/uuid"
)

type WalletEntity struct {
	Id           int       `db:"id"`
	UserPublicId uuid.UUID `db:"user_public_id"`
	Balance      int64     `db:"balance"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewWallet(userId uuid.UUID) WalletEntity {
	return WalletEntity{
		UserPublicId: userId,
		Balance:      0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}