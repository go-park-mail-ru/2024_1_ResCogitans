package csrf

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

type CSRFStorage struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.Mutex
}

func NewCSRFStorage(client *redis.Client) *CSRFStorage {
	ctx, cancel := context.WithCancel(context.Background())
	// Обеспечение освобождения ресурсов контекста при завершении работы
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return &CSRFStorage{
		client: client,
		ctx:    ctx,
		mu:     sync.Mutex{},
	}
}

func (cs *CSRFStorage) SaveToken(token string, key []byte, userID int) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	err := cs.client.Set(cs.ctx, string(key), token, 0).Err()
	if err != nil {
		return fmt.Errorf("error saving token to Redis: %w", err)
	}

	err = cs.client.Set(cs.ctx, strconv.Itoa(userID), key, 0).Err()
	if err != nil {
		return fmt.Errorf("error saving key to Redis: %w", err)
	}
	return nil
}

func (cs *CSRFStorage) GetKey(userID int) ([]byte, error) {
	return cs.client.Get(cs.ctx, strconv.Itoa(userID)).Bytes()
}

func (cs *CSRFStorage) GetToken(key []byte) (string, error) {
	return cs.client.Get(cs.ctx, string(key)).Result()
}

func (cs *CSRFStorage) DeleteToken(key []byte) error {
	return cs.client.Del(cs.ctx, string(key)).Err()
}
