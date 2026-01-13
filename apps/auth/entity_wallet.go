package auth

import (
	"time"

	"github.com/google/uuid"
)

type WalletEntity struct {
	Id           			int       `db:"id"`
	UserPublicId 			uuid.UUID `db:"user_public_id"`
	WalletPublicId 			uuid.UUID `db:"wallet_public_id"`
	Name 					string `db:"name"`
	Balance      			int64     `db:"balance"`
	CreatedAt    			time.Time `db:"created_at"`
	UpdatedAt    			time.Time `db:"updated_at"`
}

func NewWallet(userId uuid.UUID, name string) WalletEntity {
	return WalletEntity{
		UserPublicId: userId,
		Name: name,
		WalletPublicId: uuid.New(),
		Balance:      0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}