package main

import (
	"ariskaAdi/e-wallet/apps/auth"
	"ariskaAdi/e-wallet/eksternal/database"
	"ariskaAdi/e-wallet/internal/config"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	filename := "../../cmd/api/config.yaml"
	if err := config.LoadConfig(filename); err != nil {
		panic(err)
	}

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