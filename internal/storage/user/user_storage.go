package user

import (
	"sync"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/pkg/errors"
)

type StorageInterface interface {
	SaveUser(username string, password string, salt string) entities.User
	ChangeUsername(userID int, username string) error
	ChangePassword(userID int, newPassword, salt string) error
	GetUserByID(userID int) (entities.User, error)
	GetUserByUsername(username string) (entities.User, error)
	DeleteUser(userID int) error
	IsUsernameTaken(username string) bool
}

type UserStorage struct {
	Users map[int]entities.User
	mu    sync.Mutex
}

func NewUserStorage() StorageInterface {
	return &UserStorage{
		Users: make(map[int]entities.User),
		mu:    sync.Mutex{},
	}
}

func (us *UserStorage) SaveUser(username string, hashPassword string, salt string) entities.User {
	us.mu.Lock()
	defer us.mu.Unlock()

	userID := 1
	if len(us.Users) != 0 {
		userID = len(us.Users) + 1
	}

	newUser := entities.User{
		ID:       userID,
		Username: username,
		Password: hashPassword,
		Salt:     salt,
	}

	us.Users[userID] = newUser
	return newUser
}

func (us *UserStorage) ChangeUsername(userID int, newUsername string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, exists := us.Users[userID]
	if !exists {
		return errors.New("User not found")
	}
	user.Username = newUsername
	us.Users[userID] = user
	return nil
}

func (us *UserStorage) ChangePassword(userID int, newPassword, salt string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, exists := us.Users[userID]
	if !exists {
		return errors.New("User not found")
	}
	user.Password = newPassword
	user.Salt = salt
	us.Users[userID] = user
	return nil
}

func (us *UserStorage) GetUserByID(userID int) (entities.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, exists := us.Users[userID]
	if !exists {
		return entities.User{}, errors.New("User not found")
	}
	return user, nil
}

func (us *UserStorage) GetUserByUsername(username string) (entities.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	for _, user := range us.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return entities.User{}, errors.New("User not found")
}

func (us *UserStorage) DeleteUser(userID int) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.Users[userID]; !exists {
		return errors.New("User not found")
	}

	delete(us.Users, userID)
	return nil
}

func (us *UserStorage) IsUsernameTaken(username string) bool {
	us.mu.Lock()
	defer us.mu.Unlock()

	for _, user := range us.Users {
		if user.Username == username {
			return true
		}
	}
	return false
}
