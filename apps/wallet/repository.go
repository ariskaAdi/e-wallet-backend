package wallet

import (
	"ariskaAdi/e-wallet/infra/response"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) repository {
	return repository{db: db}
}

func (r repository) GetWalletByUserPublicId(ctx context.Context, userPublicId string) (model WalletEntity, err error) {
	query := `
		SELECT id, user_public_id, balance, created_at, updated_at
		FROM wallet
		WHERE user_public_id = $1
	`

	err = r.db.GetContext(ctx, &model, query, userPublicId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}

	return 
}

func (r repository) GetByWalletId(ctx context.Context, walletPublicId string) (model WalletEntity, err error) {
	query := `
		SELECT id, user_public_id, name, balance, created_at, updated_at
		FROM wallet
		WHERE wallet_public_id = $1
	`

	err = r.db.GetContext(ctx, &model, query, walletPublicId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}
	return
}

