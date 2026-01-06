package database

import (
	"ariskaAdi/e-wallet/internal/config"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg config.DBConfig) (db *sqlx.DB, err error ){
	dsn := fmt.Sprintf("host=%s port=%s  user=%s  dbname=%s password=%s sslmode=disable ",
	cfg.Host,
	cfg.Port,
	cfg.User,
	cfg.Name,
	cfg.Password,
)
	db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	db.SetMaxIdleConns(int(cfg.ConnectionPool.MaxIdle))
	db.SetMaxOpenConns(int(cfg.ConnectionPool.MaxOpen))
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnectionPool.MaxIdleTime) * time.Second) 
	db.SetConnMaxLifetime(time.Duration(cfg.ConnectionPool.MaxLifetime) * time.Second)
	
	
	return
}