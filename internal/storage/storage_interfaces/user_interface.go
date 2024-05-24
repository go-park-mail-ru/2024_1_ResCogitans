package storage

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type UserStorageInterface interface {
	SaveUser(ctx context.Context, email string, password, salt string) error
	ChangeEmail(ctx context.Context, userID int, email string) error
	ChangePassword(ctx context.Context, userID int, newPassword, salt string) error
	GetUserByID(ctx context.Context, userID int) (entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
	DeleteUser(ctx context.Context, userID int) error
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}
