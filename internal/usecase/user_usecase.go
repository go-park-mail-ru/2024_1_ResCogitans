package usecase

import (
	"net/http"
	"regexp"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseInterface interface {
	CreateUser(email string, password string) error
	GetUserByEmail(email string) (entities.User, error)
	GetUserByID(userID int) (entities.User, error)
	DeleteUser(userID int) error
	UserDataVerification(email, password string) error
	ChangeData(userID int, email, password string) (entities.User, error)
	UserExists(email, password string) error
	IsEmailTaken(email string) (bool, error)
}

type UserUseCase struct {
	UserStorage storage.UserStorageInterface
}

func NewUserUseCase(storage storage.UserStorageInterface) UserUseCaseInterface {
	return &UserUseCase{
		UserStorage: storage,
	}
}

func (u *UserUseCase) CreateUser(email, password string) error {
	salt := uuid.New().String()
	saltedPassword := password + salt
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed creating hash password")
	}

	err = u.UserStorage.SaveUser(email, string(hashPassword), salt)
	return err
}

func (u *UserUseCase) GetUserByEmail(email string) (entities.User, error) {
	return u.UserStorage.GetUserByEmail(email)
}

func (u *UserUseCase) GetUserByID(userID int) (entities.User, error) {
	return u.UserStorage.GetUserByID(userID)
}

func (u *UserUseCase) DeleteUser(userID int) error {
	return u.UserStorage.DeleteUser(userID)
}

func (u *UserUseCase) UserDataVerification(email, password string) error {
	if !ValidateEmail(email) {
		return httperrors.NewHttpError(http.StatusBadRequest, "Email doesn't meet requirements")
	}
	if !ValidatePassword(password) {
		return httperrors.NewHttpError(http.StatusBadRequest, "Password doesn't meet requirements")
	}

	return nil
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`
	matched, _ := regexp.MatchString(pattern, email)
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

func (u *UserUseCase) ChangeData(userID int, email, password string) (entities.User, error) {
	user, err := u.UserStorage.GetUserByID(userID)
	if err != nil {
		return entities.User{}, err
	}

	if user.Username != email {
		err := u.UserStorage.ChangeEmail(userID, email)
		if err != nil {
			return entities.User{}, err
		}
	}

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
	return u.UserStorage.GetUserByID(userID)
}

func (u *UserUseCase) UserExists(email, password string) error {
	user, err := u.UserStorage.GetUserByEmail(email)
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

func (u *UserUseCase) IsEmailTaken(email string) (bool, error) {
	return u.UserStorage.IsEmailTaken(email)
}
