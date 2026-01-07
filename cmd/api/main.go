package main

import (
	"ariskaAdi/e-wallet/apps/auth"
	"ariskaAdi/e-wallet/eksternal/database"
	"ariskaAdi/e-wallet/internal/config"
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

	router := fiber.New(fiber.Config{
		Prefork: true,
		AppName: config.Cfg.App.Name,
	})

	auth.Init(router, db)

	router.Listen(":" + config.Cfg.App.Port)
}