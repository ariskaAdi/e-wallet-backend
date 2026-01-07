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

func (r repository) CreateAuth(ctx context.Context, model AuthEntity) (err error) {
	query := `
		INSERT INTO auth (
			username, email, password, created_at, updated_at, public_id, otp, verified
		) VALUES (
			:username, :email, :password, :created_at, :updated_at, :public_id, :otp, :verified
		)
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, model)

	return
}

func (r repository) GetAuthByEmail(ctx context.Context, email string) (model AuthEntity, err error) {
		query := `
		SELECT id, username, email, password,  created_at, updated_at, public_id, otp, verified
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

	_, err = r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return err
	}

	return nil
}