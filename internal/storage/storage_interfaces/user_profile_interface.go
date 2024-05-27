package storage

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type UserProfileStorageInterface interface {
	GetUserProfileByID(ctx context.Context, userID int) (entities.UserProfile, error)
	EditUsername(ctx context.Context, userID int, username string) error
	EditUserBio(ctx context.Context, userID int, bio string) error
	EditUserAvatar(ctx context.Context, userID int, avatar string) error
}
