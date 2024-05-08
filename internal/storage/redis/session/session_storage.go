package session

import (
	"context"
	"strconv"
	"sync"
	"time"

	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.Mutex
}

func NewSessionStorage(client *redis.Client) storage.SessionStorageInterface {
	ctx, cancel := context.WithCancel(context.Background())
	// Обеспечение освобождения ресурсов контекста при завершении работы
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return &RedisStorage{
		client: client,
		ctx:    ctx,
		mu:     sync.Mutex{},
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
