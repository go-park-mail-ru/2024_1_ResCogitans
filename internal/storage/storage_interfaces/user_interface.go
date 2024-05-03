package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type UserStorageInterface interface {
	SaveUser(email string, password, salt string) error
	ChangeEmail(userID int, email string) error
	ChangePassword(userID int, newPassword, salt string) error
	GetUserByID(userID int) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	DeleteUser(userID int) error
	IsEmailTaken(email string) (bool, error)
}
