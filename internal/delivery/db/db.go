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
	db_info := config.DSN

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db_info.Host, db_info.Port, db_info.User, db_info.Password, db_info.DBname)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Logger().Error("Cannot open database")
		return nil, err
	}
	poolConfig.MaxConns = 10
	poolConfig.MaxConnLifetime = time.Hour

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func GetRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisData.Addr,
		Username: config.RedisData.Username,
		Password: config.RedisData.Password,
		DB:       config.RedisData.DB,
	})

	// Проверяем соединение с Redis
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
