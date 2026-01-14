package transaction

import (
	"ariskaAdi/e-wallet/infra/response"
	"time"
)

type WalletEntity struct {
	Id             int    `db:"id"`
	UserPublicId   string `db:"user_public_id"`
	WalletPublicId string `db:"wallet_public_id"`
	Name           string `db:"name"`
	Balance        int64  `db:"balance"`
	CreatedAt      time.Time `db:"created_at"`
    UpdatedAt      time.Time `db:"updated_at"`
}

func (w WalletEntity) isExist() bool {
	return w.Id != 0
}

func (w *WalletEntity) UpdateBalanceDebit(amount int64) (err error) {
	if w.Balance < amount {
		err = response.ErrInsufficientBalance
		return
	}
	w.Balance -= amount
	return
}

func (w *WalletEntity) UpdateBalanceCredit(amount int64) (err error) {
	w.Balance += amount
	return
}