package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("ENCRYPTION_SALT", 10)
	
	// default values
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("ENCRYPTION_SALT", 10)

	if err := viper.ReadInConfig(); err != nil {
		log.Println("config file not found, using env only")
	}
	
	Cfg = Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Port: viper.GetString("APP_PORT"),
			Encryption: EncryptionConfig{
				Salt: uint8(viper.GetInt("ENCRYPTION_SALT")),
				JWTSecret: viper.GetString("JWT_SECRET"),
			},
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Name:     viper.GetString("DB_NAME"),
			Password: viper.GetString("DB_PASSWORD"),
			ConnectionPool: DBConnectionPoolConfig{
				MaxIdle:     viper.GetInt("DB_MAX_IDLE"),
				MaxOpen:     viper.GetInt("DB_MAX_OPEN"),
				MaxLifetime: viper.GetInt("DB_MAX_LIFETIME"),
				MaxIdleTime: viper.GetInt("DB_MAX_IDLE_TIME"),
			},
		},
	}

	validate()
}

func validate() {
	if Cfg.App.Encryption.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if Cfg.DB.Host == "" {
		log.Fatal("DB_HOST is required")
	}
}