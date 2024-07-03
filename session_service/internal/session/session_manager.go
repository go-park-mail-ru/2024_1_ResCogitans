package session

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/gen"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/internal/storage"
	"github.com/google/uuid"
)

type SessionManager struct {
	gen.UnimplementedSessionServiceServer
	storage *storage.RedisStorage
}

func NewSessionManager(storage *storage.RedisStorage) *SessionManager {
	return &SessionManager{
		storage: storage,
	}
}

func (sm *SessionManager) CreateSession(ctx context.Context, req *gen.SaveSessionRequest) (*gen.SaveSessionResponse, error) {
	sessionID := uuid.New().String()
	err := sm.storage.SaveSession(ctx, sessionID, int(req.UserID))
	if err != nil {
		return nil, err
	}
	return &gen.SaveSessionResponse{SessionID: sessionID}, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, req *gen.GetSessionRequest) (*gen.GetSessionResponse, error) {
	userID, err := sm.storage.GetSession(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}
	return &gen.GetSessionResponse{UserID: int32(userID)}, nil
}

func (sm *SessionManager) DeleteSession(ctx context.Context, req *gen.DeleteSessionRequest) (*gen.DeleteSessionResponse, error) {
	err := sm.storage.DeleteSession(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}
	return &gen.DeleteSessionResponse{}, nil
}
