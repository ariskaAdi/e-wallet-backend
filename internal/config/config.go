package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
	Encryption EncryptionConfig `yaml:"encryption"`
}

type EncryptionConfig struct {
	Salt uint8 `yaml:"salt"`
	JWTSecret string `yaml:"jwt_secret"`
}

type DBConfig struct {
	Host           string                 `yaml:"host"`
	Port           string                 `yaml:"port"`
	User           string                 `yaml:"user"`
	Name           string                 `yaml:"name"`
	Password       string                 `yaml:"password"`
	ConnectionPool DBConnectionPoolConfig `yaml:"connection_pool"`
}

type DBConnectionPoolConfig struct {
	MaxIdle     int `yaml:"max_idle"`
	MaxOpen     int `yaml:"max_open"`
	MaxLifetime int `yaml:"max_lifetime"`
	MaxIdleTime int `yaml:"max_idle_time"`
}

var Cfg Config

func LoadConfig(filename string) (err error) {
	configByte, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	return yaml.Unmarshal(configByte, &Cfg)
}