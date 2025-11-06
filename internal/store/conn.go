package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type DBConfig struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST string
	DB_PORT string
	DB_NAME string
	DB_MaxOpenConn int
	DB_MaxIdleConn int
	DB_MaxConnLifetime time.Duration
	DB_MaxConnIdletime time.Duration
}

func ConnectDB(cfg DBConfig) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOST,
		cfg.DB_PORT, cfg.DB_NAME,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.SetMaxOpenConns(cfg.DB_MaxOpenConn)
	db.SetMaxIdleConns(cfg.DB_MaxIdleConn)
	db.SetConnMaxLifetime(cfg.DB_MaxConnLifetime)
	db.SetConnMaxIdleTime(cfg.DB_MaxConnIdletime)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed connect to database: %v", err)
	}

	log.Println("Connected to database!!")
	return db
}

func TestRedis(ctx context.Context,rd *redis.Client) {
	s, err := rd.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	log.Printf("Conneted to redis!! %s", s)
}