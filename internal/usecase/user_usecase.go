package usecase

import (
	"net/http"
	"regexp"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/user"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	CreateUser(username string, password string) (entities.User, error)
	GetUserByUsername(username string) (entities.User, error)
	GetUserByID(userID int) (entities.User, error)
	DeleteUser(userID int) error
	UserDataVerification(username, password string) error
	ChangeData(userID int, username, password string) (entities.User, error)
	UserExists(username, password string) error
	IsUsernameTaken(username string) bool
}

type UserUseCase struct {
	UserStorage user.StorageInterface
}

func NewUserUseCase(storage user.StorageInterface) UserInterface {
	return &UserUseCase{
		UserStorage: storage,
	}
}

func (u *UserUseCase) CreateUser(username, password string) (entities.User, error) {
	salt := uuid.New().String()
	saltedPassword := password + salt
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, errors.New("Failed creating hash password")
	}
	newUser := u.UserStorage.SaveUser(username, string(hashPassword), salt)
	return newUser, nil
}

func (u *UserUseCase) GetUserByUsername(username string) (entities.User, error) {
	return u.UserStorage.GetUserByUsername(username)
}

func (u *UserUseCase) GetUserByID(userID int) (entities.User, error) {
	return u.UserStorage.GetUserByID(userID)
}

func (u *UserUseCase) DeleteUser(userID int) error {
	return u.UserStorage.DeleteUser(userID)
}

func (u *UserUseCase) UserDataVerification(username, password string) error {
	if !ValidateUsername(username) {
		return httperrors.NewHttpError(http.StatusBadRequest, "Username doesn't meet requirements")
	}
	if !ValidatePassword(password) {
		return httperrors.NewHttpError(http.StatusBadRequest, "Password doesn't meet requirements")
	}

	return nil
}

func ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]{2,}$`, username)
	return matched
}

func ValidatePassword(password string) bool {
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasMinLength := len(password) >= 8
	hasMaxLength := len(password) <= 40
	return hasDigit && hasUppercase && hasLowercase && hasMinLength && hasMaxLength
}

func (u *UserUseCase) ChangeData(userID int, username, password string) (entities.User, error) {
	user, err := u.UserStorage.GetUserByID(userID)
	if err != nil {
		return entities.User{}, err
	}

	if user.Username != username {
		err := u.UserStorage.ChangeUsername(userID, username)
		if err != nil {
			return entities.User{}, err
		}
	}
	user.Username = username

	salt := uuid.New().String()
	saltedPassword := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, errors.Wrap(err, "failed creating hash password")
	}
	err = u.UserStorage.ChangePassword(userID, string(hashedPassword), salt)
	if err != nil {
		return entities.User{}, err
	}
	user.Password = string(hashedPassword)
	user.Salt = salt
	return user, nil
}

func (u *UserUseCase) UserExists(username, password string) error {
	user, err := u.UserStorage.GetUserByUsername(username)
	if err != nil {
		return err
	}

	saltedPassword := password + user.Salt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(saltedPassword))
	if err != nil {
		return errors.Wrap(err, "password is incorrect")
	}

	return nil
}

func (u *UserUseCase) IsUsernameTaken(username string) bool {
	return u.UserStorage.IsUsernameTaken(username)
}
