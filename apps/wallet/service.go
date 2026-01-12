package wallet

import (
	"ariskaAdi/e-wallet/infra/response"
	"context"
)

type Repository interface {
	GetWalletByUserPublicId(ctx context.Context, userPublicId string) (model WalletEntity, err error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) service {
	return service{repo: repo}
}


func (s service) GetMyWallet(ctx context.Context, userPublicId string) (myWallet WalletEntity, err error) {
	myWallet, err = s.repo.GetWalletByUserPublicId(ctx, userPublicId)

	if err != nil {
		if err == response.ErrNotFound {
			myWallet = WalletEntity{}
			return myWallet, nil
		}
		return
	}

	return
}

