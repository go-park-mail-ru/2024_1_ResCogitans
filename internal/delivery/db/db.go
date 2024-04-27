package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetPostgress gets connection to postgres database
func GetPostgres() (*pgxpool.Pool, error) {
	log := logger.Logger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load database config", "error", err)
		return nil, err
	}

	dsn := buildDSN(cfg.Dsn)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error("Failed to parse database config", "error", err)
		return nil, err
	}

	poolConfig.MaxConns = 10
	poolConfig.MaxConnLifetime = time.Hour

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	return pool, nil
}

func buildDSN(cfg config.Dsn) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
}

func GetRedis() (*redis.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Проверяем соединение с Redis
	_, err = rdb.Ping(rdb.Context()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
