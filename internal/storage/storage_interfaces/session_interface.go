package storage

import (
	"context"
)

type SessionStorageInterface interface {
	SaveSession(ctx context.Context, sessionID string, userID int) error
	GetSession(ctx context.Context, sessionID string) (int, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
