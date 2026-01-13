package wallet

import (
	"time"

	"github.com/google/uuid"
)

type WalletEntity struct {
	Id          	 	int       `db:"id"`
	WalletPublicId 		uuid.UUID `db:"wallet_public_id"`
	UserPublicId 		string 	`db:"user_public_id"`
	Name 				string `db:"name"`
	Balance      		int64     `db:"balance"`
	CreatedAt    		time.Time `db:"created_at"`
	UpdatedAt    		time.Time `db:"updated_at"`
}



