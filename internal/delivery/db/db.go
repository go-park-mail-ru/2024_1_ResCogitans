package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// GetPostgres соединение с базой данных
func GetPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := buildDSN(cfg.Dsn)

	// создание конфигурации базы данных
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse database config")
	}

	poolConfig.MaxConns = 10               // Максимальное число пользователей
	poolConfig.MaxConnLifetime = time.Hour // Время жизни соединения

	// Соединение с базой данных с использованием созданной конфигурацией
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return pool, nil
}

// buildDSN создание пути для базы данных
func buildDSN(cfg config.Dsn) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
}

func GetCSRFRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.CSRF.Host, cfg.CSRF.Port),
		Password: cfg.CSRF.Password,
		DB:       cfg.CSRF.DB,
	})

	// Проверка соединения с базой данных
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
