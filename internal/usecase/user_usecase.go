package usecase

import (
	"context"
	"net/http"
	"regexp"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/user"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseInterface interface {
	CreateUser(ctx context.Context, email string, password string) error
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
	GetUserByID(ctx context.Context, userID int) (entities.User, error)
	DeleteUser(ctx context.Context, userID int) error
	UserDataVerification(email, password string) error
	ChangePassword(ctx context.Context, userID int, password string) (entities.User, error)
	UserExists(ctx context.Context, email, password string) error
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}

type UserUseCase struct {
	UserStorage *user.UserStorage
}

func NewUserUseCase(storage *user.UserStorage) *UserUseCase {
	return &UserUseCase{
		UserStorage: storage,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, email, password string) error {
	salt := uuid.New().String()
	saltedPassword := password + salt
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed creating hash password")
	}

	err = u.UserStorage.SaveUser(ctx, email, string(hashPassword), salt)
	return err
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	return u.UserStorage.GetUserByEmail(ctx, email)
}

func (u *UserUseCase) GetUserByID(ctx context.Context, userID int) (entities.User, error) {
	return u.UserStorage.GetUserByID(ctx, userID)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, userID int) error {
	return u.UserStorage.DeleteUser(ctx, userID)
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

func (u *UserUseCase) ChangePassword(ctx context.Context, userID int, password string) (entities.User, error) {
	salt := uuid.New().String()
	saltedPassword := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, errors.Wrap(err, "failed creating hash password")
	}
	err = u.UserStorage.ChangePassword(ctx, userID, string(hashedPassword), salt)
	if err != nil {
		return entities.User{}, err
	}
	return u.UserStorage.GetUserByID(ctx, userID)
}

func (u *UserUseCase) UserExists(ctx context.Context, email, password string) error {
	user, err := u.UserStorage.GetUserByEmail(ctx, email)
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

func (u *UserUseCase) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	return u.UserStorage.IsEmailTaken(ctx, email)
}
