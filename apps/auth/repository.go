package auth

import (
	"ariskaAdi/e-wallet/infra/response"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) repository {
	return repository{db: db}
}

func (r repository) CreateAuth(ctx context.Context, tx *sqlx.Tx, model AuthEntity) (err error) {
	query := `
		INSERT INTO auth (
			username, email, password, created_at, updated_at, public_id, otp, verified
		) VALUES (
			:username, :email, :password, :created_at, :updated_at, :public_id, :otp, :verified
		)
	`

	_, err = tx.NamedExecContext(ctx, query, model)
	if err != nil {
		return
	}

	return
}

func (r repository) GetAuthByEmail(ctx context.Context, email string) (model AuthEntity, err error) {
		query := `
		SELECT id, username, email, password,  created_at, updated_at, public_id, verified
		FROM auth
		WHERE email = $1
	`
	err = r.db.GetContext(ctx, &model, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}
	return
}

func (r repository) UpdateAuthVerifiedOtp(ctx context.Context, model AuthEntity) (err error) {

	model.Verified = true
	model.UpdatedAt = time.Now()

	query := `
		UPDATE auth
		SET verified = :verified, otp = null, updated_at = :updated_at
		WHERE email = :email AND otp = :otp
	`

	res, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return response.ErrOtpInvalid
	}

	return nil
}

// Begin implements Repository.
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

func (r repository) CreateWallet(ctx context.Context, tx *sqlx.Tx, model WalletEntity) (err error) {
	query := `
		INSERT INTO wallet (
			user_public_id, balance, created_at, updated_at
		) VALUES (
			:user_public_id, :balance, :created_at, :updated_at
		)
	`
	_, err = r.db.NamedExecContext(ctx, query, model)
	return
}

