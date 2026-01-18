package topup

import (
	infrafiber "ariskaAdi/e-wallet/infra/fiber"
	"ariskaAdi/e-wallet/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(router fiber.Router, db *sqlx.DB, cfg *config.Config) {
	// REPOSITORY 
	repo := newRepository(db)

	// PAYMENT GATEWAY
	payment := newXenditService(cfg)

	// SERVICE
	svc := newService(repo, repo, repo, payment)

	handler := newHandler(svc)

	topupRoute := router.Group("/topup")
	{
		topupRoute.Use(infrafiber.CheckAuth())
		topupRoute.Post("/create", handler.CreateTopUp)
	}

	router.Post("/xendit/callback", handler.XenditCallback)

}