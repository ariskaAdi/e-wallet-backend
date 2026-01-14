package transaction

import (
	"ariskaAdi/e-wallet/infra/response"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	TransactionDBRepository
	TransactionRepository
	WalletRepository
	InquiryRepository
}

type TransactionDBRepository interface {
	Begin(ctx context.Context) (tx *sqlx.Tx, err error)
	Rollback(ctx context.Context, tx *sqlx.Tx) (err error)
	Commit(ctx context.Context, tx *sqlx.Tx) (err error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *sqlx.Tx, model TransactionEntity) (err error)
}

type WalletRepository interface {
	GetByUserPublicId(ctx context.Context, userPublicId string) (model WalletEntity, err error)
	GetByWalletPublicId(ctx context.Context, walletPublicId string) (model WalletEntity, err error)
	GetByUserPublicIdForUpdate(ctx context.Context, tx *sqlx.Tx, userPublicId string) (model WalletEntity, err error)
	GetByWalletPublicIdForUpdate(ctx context.Context, tx *sqlx.Tx, walletPublicId string) (model WalletEntity, err error)
	UpdateWallet(ctx context.Context, tx *sqlx.Tx, model WalletEntity) (err error)
}

type InquiryRepository interface {
	CreateInquiry(ctx context.Context, model InquiryEntity) (err error)
	GetInquiryByKey(ctx context.Context, inquiryKey string) (model InquiryEntity, err error)
	DeleteInquiry(ctx context.Context, inquiryKey string) (err error)
}


type service struct {
	repo Repository
}

func newService(repo Repository) service {
	return service{repo: repo}
}

func (s service) TransferInquiry(ctx context.Context, req TransferInquiryRequestPayload, UserPublicId string) (InquiryKey InquiryEntity, destWallet WalletEntity, err error) {

	// check wallet sender
	myWallet, err := s.repo.GetByUserPublicId(ctx, UserPublicId)
	if err != nil || !myWallet.isExist() {
		return InquiryEntity{}, WalletEntity{}, response.ErrNotFound
	}

		// validate
	if myWallet.WalletPublicId == destWallet.WalletPublicId {
		err = response.ErrSameWallet
		return
	}
	

	// check wallet receiver
	destWallet, err = s.repo.GetByWalletPublicId(ctx, req.Dof)
	if err != nil || !destWallet.isExist() {
		return InquiryEntity{}, WalletEntity{}, response.ErrNotFound
	}

	inquiry := NewInquiry(req, UserPublicId)
	if err := s.repo.CreateInquiry(ctx, inquiry); err != nil {
		return InquiryEntity{}, WalletEntity{}, nil
	}

	return inquiry, destWallet, nil
}

func (s service) TransferExecute(
	ctx context.Context,
	req TransferExecuteRequest,
	userPublicId string,
) (err error) {

	// GET INQUIRY
	inquiry, err := s.repo.GetInquiryByKey(ctx, req.InquiryKey)
	if err != nil {
		return response.ErrInquiryNotFound
	}

	if time.Now().After(inquiry.ExpiredAt){
		return response.ErrInquiryExpired
	}

	if inquiry.Sof != userPublicId {
		return response.ErrUnauthorized
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = s.repo.Rollback(ctx, tx)
		}
	}()

	// source wallet (LOCK)
	sourceWallet, err := s.repo.GetByUserPublicIdForUpdate(ctx, tx, userPublicId)
	if err != nil {
		return
	}

	// destination wallet (LOCK)
	destWallet, err := s.repo.GetByWalletPublicIdForUpdate(
		ctx,
		tx,
		inquiry.Dof,
	)
	if err != nil {
		return
	}

	if err = sourceWallet.UpdateBalanceDebit(req.Amount); err != nil {
		return
	}
	_ = destWallet.UpdateBalanceCredit(req.Amount)

	debitTx := NewDebitTransaction(
		sourceWallet,
		destWallet,
		req.Amount,
		req.Description,
	)

	creditTx := NewCreditTransaction(
		sourceWallet,
		destWallet,
		req.Amount,
		req.Description,
	)

	if err = s.repo.CreateTransaction(ctx, tx, debitTx); err != nil {
		return
	}
	if err = s.repo.CreateTransaction(ctx, tx, creditTx); err != nil {
		return
	}

	if err = s.repo.UpdateWallet(ctx, tx, sourceWallet); err != nil {
		return
	}
	if err = s.repo.UpdateWallet(ctx, tx, destWallet); err != nil {
		return
	}

	if err = s.repo.Commit(ctx, tx); err != nil {
		return
	}
	_ = s.repo.DeleteInquiry(ctx, req.InquiryKey)
	return
}
