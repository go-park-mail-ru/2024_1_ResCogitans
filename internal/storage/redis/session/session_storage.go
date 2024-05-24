package session

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisStorage struct {
	client *redis.Client
}

func NewSessionStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

func (rs *RedisStorage) SaveSession(ctx context.Context, sessionID string, userID int) error {
	err := rs.client.Set(ctx, sessionID, userID, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisStorage) GetSession(ctx context.Context, sessionID string) (int, error) {
	userIDStr, err := rs.client.Get(ctx, sessionID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, errors.New("Session not found")
		}
		return 0, err
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (rs *RedisStorage) DeleteSession(ctx context.Context, sessionID string) error {
	err := rs.client.Del(ctx, sessionID).Err()
	if err != nil {
		return err
	}
	return nil
}
