package database

import (
	"fmt"

	cfg "github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/config"
	"github.com/go-redis/redis/v8"
)

func GetSessionRedis() (*redis.Client, error) {
	config, err := cfg.LoadConfig()
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// Проверяем соединение с Redis
	_, err = rdb.Ping(rdb.Context()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
