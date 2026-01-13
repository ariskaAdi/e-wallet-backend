package auth

import (
	"ariskaAdi/e-wallet/infra/response"
	"ariskaAdi/e-wallet/internal/config"
	"ariskaAdi/e-wallet/internal/mail"
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	AuthRepository
	WalletRepository
	AuthDBRepository
}


type AuthRepository interface {
	GetAuthByEmail(ctx context.Context, email string) (model AuthEntity, err error)
	CreateAuth(ctx context.Context, tx *sqlx.Tx, model AuthEntity) (err error)
	UpdateAuthVerifiedOtp(ctx context.Context, model AuthEntity) (err error)
}

type AuthDBRepository interface {
	Begin(ctx context.Context) (tx *sqlx.Tx, err error)
	Rollback(ctx context.Context, tx *sqlx.Tx,) ( err error)
	Commit(ctx context.Context, tx *sqlx.Tx,) ( err error)
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, tx *sqlx.Tx, model WalletEntity) (err error)
}

type service struct {
	repo        Repository
	emailWorker *mail.Worker
}

func newService(repo Repository, emailWorker *mail.Worker) service {
	return service{
		repo: repo,
		emailWorker: emailWorker,
	}
}

func (s service) register(ctx context.Context, req RegisterRequestPayload) (err error) {
	authEntity := NewFormRegisterRequset(req)

	if err = authEntity.Validate(); err != nil {
		return
	}

	if err = authEntity.EncryptPassword(int(config.Cfg.App.Encryption.Salt)); err != nil {
		return
	}

	

	model, err := s.repo.GetAuthByEmail(ctx, authEntity.Email)
	if err != nil {
		if err != response.ErrNotFound {
			return
		}
	}

	if model.IsExist() {
		return response.ErrEmailAlredyExist
	}

	// start transaction
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return
	}

	defer s.repo.Rollback(ctx, tx)

	// create user
	if err = s.repo.CreateAuth(ctx, tx, authEntity); err != nil {
		return	
	}

	// create wallet
	wallet := NewWallet(authEntity.PublicId, authEntity.Username)
	if err = s.repo.CreateWallet(ctx, tx, wallet); err != nil {
		return
	}

	if err = s.repo.Commit(ctx, tx); err != nil {
		return
	}

	s.emailWorker.Enqueue(mail.OTPJob{
		To:       authEntity.Email,
		Username: authEntity.Username,
		Otp:      authEntity.OTP,
	})
	return nil
}

func (s service) login(ctx context.Context, req LoginRequestPayload) (token string, err error) {
	authEntity := NewFormLoginRequset(req)

	if err = authEntity.EmailValidate(); err != nil {
		return
	}

	if err = authEntity.PasswordValidate(); err != nil {
		return
	}

	model, err := s.repo.GetAuthByEmail(ctx, authEntity.Email)
	if err != nil { 
		return	
	}

	if err = authEntity.VerifyPasswordFromPlain(model.Password); err != nil {
		err = response.ErrPasswordNotMatch
		return
	}

	token, err = model.GenerateToken(config.Cfg.App.Encryption.JWTSecret)
	return
}

func (s service) verifyOtp(ctx context.Context, req ValidateOtpRequestPayload) (err error) {
	authEntity := NewFormValidateOtpRequset(req)

	if err = authEntity.EmailValidate(); err != nil {
		return
	}

	if err = authEntity.OtpValidate(); err != nil {
		return
	}

	model, err := s.repo.GetAuthByEmail(ctx, authEntity.Email)
	if err != nil {
		return err
	}

	if model.Verified == true {
		return response.ErrEmailAlreadyVerified
	}

	model.OTP = authEntity.OTP

	if err = s.repo.UpdateAuthVerifiedOtp(ctx, model); err != nil {
		return response.ErrOtpInvalid	
	}

	return nil
}