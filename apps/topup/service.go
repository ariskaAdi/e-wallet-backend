package topup

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TopUpRepository interface {
	FindTopUpByTopUpIdForUpdate(ctx context.Context, tx *sqlx.Tx, topUpId string) (model TopUpEntity, err error)
	UpdateWallet(ctx context.Context, tx *sqlx.Tx, model WalletEntity) (err error)
	Create(ctx context.Context, tx *sqlx.Tx, model TopUpEntity) (err error)
	Update(ctx context.Context, tx *sqlx.Tx, model TopUpEntity) (err error)
}

type TxManager interface {
	Begin(ctx context.Context) (tx *sqlx.Tx, err error)
	Rollback(ctx context.Context, tx *sqlx.Tx,) ( err error)
	Commit(ctx context.Context, tx *sqlx.Tx,) ( err error)
}

type PaymentGateway interface {
	GenerateSnapURL(ctx context.Context, model *TopUpEntity) (err error)
	CheckPayment(ctx context.Context, orderId string) (OrderId string, err error)
}

type Service interface {
	CreateTopUp(ctx context.Context, amount int64, userPublicId string) (model TopUpEntity, err error)
	ConfirmTopUp(ctx context.Context, topUpId string) (err error)
	FailTopUp(ctx context.Context, topUpId string) (err error)
}

type WalletRepository interface {
	AddBalance(ctx context.Context, tx *sqlx.Tx, userPublicId string, amount int64) error
}

type service struct {
	repo    TopUpRepository
	walletRepo WalletRepository
	tx      TxManager
	payment PaymentGateway
}


func newService(repo TopUpRepository, walletRepo WalletRepository, tx TxManager, payment PaymentGateway ) Service {
	return &service{
		repo:    repo,
		walletRepo: walletRepo,
		tx:      tx,
		payment: payment,
	}
}

func (s *service) CreateTopUp(ctx context.Context, amount int64, userPublicId string) (TopUpEntity, error) {

	topup := NewTopUp(amount, userPublicId)

	tx, err := s.tx.Begin(ctx)
	if err != nil {
		return TopUpEntity{}, err
	}
	defer s.tx.Rollback(ctx, tx)

	if err := s.repo.Create(ctx, tx, topup); err != nil {
		return TopUpEntity{}, err
	}

	if err := s.payment.GenerateSnapURL(ctx, &topup); err != nil {
		return TopUpEntity{}, err
	}

	if err := s.repo.Update(ctx, tx, topup); err != nil {
		return TopUpEntity{}, err
	}

	if err := s.tx.Commit(ctx, tx); err != nil {
		return TopUpEntity{}, err
	}

	return topup, nil
}


func (s *service) ConfirmTopUp(ctx context.Context, topUpId string ) error {

	tx, err := s.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.tx.Rollback(ctx, tx)

	topup, err := s.repo.FindTopUpByTopUpIdForUpdate(ctx, tx, topUpId)
	if err != nil {
		return err
	}

	if topup.Status != TopUpPending {
		return nil
	}

	if err := topup.MarkSuccess(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, tx, topup); err != nil {
		return err
	}

	if err := s.walletRepo.AddBalance(
		ctx,
		tx,
		topup.UserPublicId,
		topup.Amount,
	); err != nil {
		return err
	}

	return s.tx.Commit(ctx, tx)
}

func (s *service) FailTopUp(ctx context.Context, topupId string) error {

	tx, err := s.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.tx.Rollback(ctx, tx)

	topup, err := s.repo.FindTopUpByTopUpIdForUpdate(ctx, tx, topupId)
	if err != nil {
		return err
	}

	if topup.Status != TopUpPending {
		return nil
	}

	topup.MarkFailed()

	if err := s.repo.Update(ctx, tx, topup); err != nil {
		return err
	}

	return s.tx.Commit(ctx, tx)
}
