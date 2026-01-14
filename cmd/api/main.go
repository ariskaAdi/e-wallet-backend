package main

import (
	"ariskaAdi/e-wallet/apps/auth"
	"ariskaAdi/e-wallet/apps/transaction"
	"ariskaAdi/e-wallet/apps/wallet"
	"ariskaAdi/e-wallet/eksternal/database"
	"ariskaAdi/e-wallet/internal/config"
	"ariskaAdi/e-wallet/internal/mail"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
	log.Println(".env not found, using OS env")
}

	config.LoadConfig()

	db, err := database.ConnectPostgres(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	if db != nil {
		log.Println("DB CONNECTED")
	}

	// worker email
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emailWorker := mail.NewWorker(mail.SMTPConfig{
		Host:     config.Cfg.SMTP.Host,
		Port:     config.Cfg.SMTP.Port,
		User:     config.Cfg.SMTP.User,
		Pass:     config.Cfg.SMTP.Pass,
		From:     config.Cfg.SMTP.From,
	}, 100)
	emailWorker.Start(ctx)

	router := fiber.New(fiber.Config{
		Prefork: true,
		AppName: config.Cfg.App.Name,
	})

	auth.Init(router, db, emailWorker)
	wallet.Init(router, db)
	transaction.Init(router, db)

	router.Listen(":" + config.Cfg.App.Port)
}