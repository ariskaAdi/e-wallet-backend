package transaction

import (
	"ariskaAdi/e-wallet/infra/response"
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	TransactionDBRepository
	TransactionRepository
	WalletRepository
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
	GetByWalletPublicIdForUpdate(ctx context.Context, tx *sqlx.Tx, req TransferExecuteRequest) (model WalletEntity, err error)


	UpdateWallet(ctx context.Context, tx *sqlx.Tx, model WalletEntity) (err error)
}


type service struct {
	repo Repository
}

func newService(repo Repository) service {
	return service{repo: repo}
}

func (s service) TransferInquiry(ctx context.Context, req TransferInquiryRequestPayload, UserPublicId string) (destWallet WalletEntity, err error) {
	
	// check wallet sender
	myWallet, err := s.repo.GetByUserPublicId(ctx, UserPublicId)
	if err != nil {
		return
	}
	if !myWallet.isExist() {
		err = response.ErrNotFound
		return
	}

	// check wallet receiver
	destWallet, err = s.repo.GetByWalletPublicId(ctx, req.DestinationWalletPublicId)
	if err != nil {
		return
	}
	if !destWallet.isExist() {
		err = response.ErrNotFound
		return
	}

	// validate
	if myWallet.WalletPublicId == destWallet.WalletPublicId {
		err = response.ErrSameWallet
		return
	}

	return
}

func (s service) TransferExecute(
	ctx context.Context,
	req TransferExecuteRequest,
	userPublicId string,
) (err error) {

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
		req.DestinationWalletPublicId, // dari inquiry cache
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

	err = s.repo.Commit(ctx, tx)
	return
}
