package topup

import (
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

func (r repository) Begin(ctx context.Context) (tx *sqlx.Tx, err error) {
	tx, err = r.db.BeginTxx(ctx, &sql.TxOptions{})
	return
}

// Commit implements Repository.
func (repository) Commit(ctx context.Context, tx *sqlx.Tx) (err error) {
	return tx.Commit()
}

// Rollback implements Repository.
func (repository) Rollback(ctx context.Context, tx *sqlx.Tx) (err error) {
	return tx.Rollback()
}

func (r repository) Create(ctx context.Context, tx *sqlx.Tx, model TopUpEntity) (err error) {
	query := `
		INSERT INTO topup (
			topup_id, user_public_id, amount, status, snap_url, created_at, updated_at
		) VALUES (
			:topup_id, :user_public_id, :amount, :status, :snap_url, :created_at, :updated_at
		)
	`
	_, err = tx.NamedExecContext(ctx, query, model)
	return
}

func (r repository) Update(ctx context.Context, tx *sqlx.Tx,  model TopUpEntity) (err error) {
	query := `
		UPDATE topup
		SET status = :status, snap_url = :snap_url, updated_at = :updated_at
		WHERE topup_id = :topup_id
	`
	_, err = tx.NamedExecContext(ctx, query, model)
	return
}

func (r repository) FindTopUpByTopUpIdForUpdate(ctx context.Context, tx *sqlx.Tx, topUpId string) (model TopUpEntity, err error) {
	query := `
		SELECT id, topup_id, user_public_id, amount, status, snap_url, created_at, updated_at
		FROM topup
		WHERE topup_id = $1
		FOR UPDATE
	`
	err = tx.GetContext(ctx, &model, query, topUpId)
	return
}

func (r repository) UpdateWallet(ctx context.Context, tx *sqlx.Tx, model WalletEntity) (err error) {
	query := `
		UPDATE wallet
		SET balance = :balance,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err = tx.NamedExecContext(ctx, query, model)
	return
}

func (r repository) AddBalance(
	ctx context.Context,
	tx *sqlx.Tx,
	userPublicId string,
	amount int64,
) error {

	query := `
		UPDATE wallet
		SET balance = balance + $1,
		updated_at = now()
		WHERE user_public_id = $2
	`

	_, err := tx.ExecContext(ctx, query, amount, userPublicId)
	return err
}
