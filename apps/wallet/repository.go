package wallet

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) repository {
	return repository{db: db}
}


func (r repository) CreateWallet(ctx context.Context, model WalletEntity) (err error) {
	query := `
		INSERT INTO wallet (
			user_id, balance, created_at, updated_at
		) VALUES (
			:user_id, :balance, :created_at, :updated_at
		)
	`
	_, err = r.db.NamedExecContext(ctx, query, model)
	return
}