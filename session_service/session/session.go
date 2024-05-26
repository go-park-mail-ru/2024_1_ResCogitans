package session

import (
	"context"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/service/gen"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type SessionManager struct {
	gen.UnimplementedSessionServiceServer
	storage *redis.Client
}

func NewSessionManager(client *redis.Client) *SessionManager {
	return &SessionManager{
		storage: client,
	}
}

func (sm *SessionManager) CreateSession(ctx context.Context, req *gen.SaveSessionRequest) (*gen.SaveSessionResponse, error) {
	err := sm.storage.Set(ctx, req.SessionID, req.UserID, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}
	println("created session")
	return &gen.SaveSessionResponse{}, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, req *gen.GetSessionRequest) (*gen.GetSessionResponse, error) {
	userIDStr, err := sm.storage.Get(ctx, req.SessionID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, err
	}
	return &gen.GetSessionResponse{UserID: int32(userID)}, nil
}

func (sm *SessionManager) DeleteSession(ctx context.Context, req *gen.DeleteSessionRequest) (*gen.DeleteSessionResponse, error) {
	err := sm.storage.Del(ctx, req.SessionID).Err()
	if err != nil {
		return nil, err
	}
	return &gen.DeleteSessionResponse{}, nil
}
