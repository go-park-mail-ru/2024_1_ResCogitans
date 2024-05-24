package csrf

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type CSRFStorage struct {
	client *redis.Client
}

func NewCSRFStorage(client *redis.Client) *CSRFStorage {
	return &CSRFStorage{
		client: client,
	}
}

func (cs *CSRFStorage) SaveToken(ctx context.Context, token string, key []byte, userID int) error {
	err := cs.client.Set(ctx, string(key), token, 0).Err()
	if err != nil {
		return fmt.Errorf("error saving token to Redis: %w", err)
	}

	err = cs.client.Set(ctx, strconv.Itoa(userID), key, 0).Err()
	if err != nil {
		return fmt.Errorf("error saving key to Redis: %w", err)
	}
	return nil
}

func (cs *CSRFStorage) GetKey(ctx context.Context, userID int) ([]byte, error) {
	return cs.client.Get(ctx, strconv.Itoa(userID)).Bytes()
}

func (cs *CSRFStorage) GetToken(ctx context.Context, key []byte) (string, error) {
	return cs.client.Get(ctx, string(key)).Result()
}

func (cs *CSRFStorage) DeleteToken(ctx context.Context, key []byte) error {
	return cs.client.Del(ctx, string(key)).Err()
}
