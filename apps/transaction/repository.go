package transaction

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

// Begin implements Repository.
func (r repository) Begin(ctx context.Context) (tx *sqlx.Tx, err error) {
	tx, err = r.db.BeginTxx(ctx, &sql.TxOptions{})
	return
}

// Commit implements Repository.
func (r repository) Commit(ctx context.Context, tx *sqlx.Tx) (err error) {
	return tx.Commit()
}

// Rollback implements Repository.
func (r repository) Rollback(ctx context.Context, tx *sqlx.Tx) (err error) {
	return tx.Rollback()
}

func (r repository) CreateTransaction(ctx context.Context, tx *sqlx.Tx, model TransactionEntity) (err error) {
	query := `
		INSERT INTO transaction (
			transaction_id, wallet_public_id, sof_number, dof_number, type, amount, status, reference, description, created_at, updated_at
		) VALUES (
			:transaction_id, :wallet_public_id, :sof_number, :dof_number, :type, :amount, :status, :reference, :description, :created_at, :updated_at
		)
	`
	_, err = tx.NamedExecContext(ctx, query, model)
	return
}

func (r repository) GetByUserPublicId(ctx context.Context, userPublicId string) (model WalletEntity, err error) {
	query := `
		SELECT id, user_public_id, name, balance, created_at, updated_at
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

func (r repository) GetByWalletPublicId(ctx context.Context, walletPublicId string) (model WalletEntity, err error) {
	query := `
		SELECT id, user_public_id, wallet_public_id, name, balance, created_at, updated_at
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

func (r repository) GetByUserPublicIdForUpdate(ctx context.Context, tx *sqlx.Tx, userPublicId string) (model WalletEntity, err error) {
	query := `
		SELECT id, wallet_public_id, user_public_id, name, balance, created_at, updated_at
		FROM wallet
		WHERE user_public_id = $1
		FOR UPDATE
	`
	err = tx.GetContext(ctx, &model, query, userPublicId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}

	return

}

func (r repository) GetByWalletPublicIdForUpdate(ctx context.Context, tx *sqlx.Tx, walletPublicId string) (model WalletEntity, err error) {
	query := `
		SELECT id, wallet_public_id, user_public_id, name, balance, created_at, updated_at
		FROM wallet
		WHERE wallet_public_id = $1
		FOR UPDATE
	`

	err = tx.GetContext(ctx, &model, query, walletPublicId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}
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

func (r repository) CreateInquiry(ctx context.Context, model InquiryEntity) (err error) {
	query := `
		INSERT INTO transfer_inquiry (
			inquiry_key, sof_number, dof_number, created_at, expired_at
		) VALUES (
			:inquiry_key, :sof_number, :dof_number, :created_at, :expired_at 
		)
	`

	_, err = r.db.NamedExecContext(ctx, query, model)
	return
}

func (r repository) GetInquiryByKey(ctx context.Context, inquiryKey string) (model InquiryEntity, err error) {
	query := `
		SELECT inquiry_key, sof_number, dof_number, created_at, expired_at
		FROM transfer_inquiry
		WHERE inquiry_key = $1
	`
	err = r.db.GetContext(ctx, &model, query, inquiryKey)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}
	return

}

func (r repository) DeleteInquiry(ctx context.Context, inquiryKey string) (err error) {
	query := `
		DELETE FROM transfer_inquiry
		WHERE inquiry_key = $1
	`
	_, err = r.db.ExecContext(ctx, query, inquiryKey)
	return
}