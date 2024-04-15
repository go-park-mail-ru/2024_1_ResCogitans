package session

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.Mutex
}

func NewRedisStorage(address string, username, password string, db int) StorageInterface {
	return &RedisStorage{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Username: username,
			Password: password,
			DB:       db,
		}),
		ctx: context.Background(),
		mu:  sync.Mutex{},
	}
}

func (rs *RedisStorage) SaveSession(sessionID string, userID int) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	err := rs.client.Set(rs.ctx, sessionID, userID, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisStorage) GetSession(sessionID string) (int, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	userIDStr, err := rs.client.Get(rs.ctx, sessionID).Result()
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

func (rs *RedisStorage) DeleteSession(sessionID string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	err := rs.client.Del(rs.ctx, sessionID).Err()
	if err != nil {
		return err
	}
	return nil
}
