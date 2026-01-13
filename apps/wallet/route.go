package wallet

import (
	infrafiber "ariskaAdi/e-wallet/infra/fiber"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(router fiber.Router, db *sqlx.DB) {
	repo := newRepository(db)
	svc := newService(repo)
	handler := newHandler(svc)

	walletRoute := router.Group("/wallet")
	{
		// MIDDLEWARE
		walletRoute.Use(infrafiber.CheckAuth())

		// wallet endpoint
		walletRoute.Get("/my-wallet", handler.GetMyWallet)
		walletRoute.Get("/someone-wallet", handler.GetSomeoneWallet)
	}
}