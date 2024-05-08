package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type UserProfileStorageInterface interface {
	GetUserProfileByID(userID int) (entities.UserProfile, error)
	EditUsername(userID int, username string) error
	EditUserBio(userID int, bio string) error
	EditUserAvatar(userID int, avatar string) error
}
