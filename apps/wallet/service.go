package wallet

import "context"

type Repository interface {
	GetWalletByUserId(ctx context.Context, userId string) (model WalletEntity, err error)
	CreateWallet(ctx context.Context, model WalletEntity) (err error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) service {
	return service{repo: repo}
}


func (s service) GetWallet(ctx context.Context, userId string) (model WalletEntity, err error) {
	model , err = s.repo.GetWalletByUserId(ctx, userId)
	if err != nil {
		return
	}
	return
}

func (s service) CreateWallet(ctx context.Context, model WalletEntity) (err error) {
	err = s.repo.CreateWallet(ctx, model)
	return
}