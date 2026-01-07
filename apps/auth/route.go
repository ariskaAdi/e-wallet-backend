package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Init(r fiber.Router, db *sqlx.DB) {
	repo := newRepository(db)
	svc := newService(repo)
	handler := newHandler(svc)
	
	authRouter := r.Group("/auth")
	{
		authRouter.Post("/register", handler.register)
		authRouter.Post("/login", handler.login)
		authRouter.Post("/verify-otp", handler.verifyOtp)
	}
}