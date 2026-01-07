package auth

import (
	"ariskaAdi/e-wallet/internal/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(r fiber.Router, db *sqlx.DB, mail *mail.Worker) {
	repo := newRepository(db)
	svc := newService(repo, mail)
	handler := newHandler(svc)
	
	authRouter := r.Group("/auth")
	{
		authRouter.Post("/register", handler.register)
		authRouter.Post("/login", handler.login)
		authRouter.Post("/verify-otp", handler.verifyOtp)
	}
}