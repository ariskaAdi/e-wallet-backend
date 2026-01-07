package config

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	Name       string
	Port       string
	Encryption EncryptionConfig
}

type EncryptionConfig struct {
	Salt      uint8
	JWTSecret string
}

type DBConfig struct {
	Host           string
	Port           string
	User           string
	Name           string
	Password       string
	ConnectionPool DBConnectionPoolConfig
}

type DBConnectionPoolConfig struct {
	MaxIdle     int
	MaxOpen     int
	MaxLifetime int
	MaxIdleTime int
}

var Cfg Config
