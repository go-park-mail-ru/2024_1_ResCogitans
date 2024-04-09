package registration

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"
	userRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type RegistrationHandler struct{}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

var (
	errCreateUser = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed creating new profile",
	}
	errSetSession = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed setting session",
	}
	errInternal = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "internal Error",
	}
)

// func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (UserResponse, error) {
// 	username := requestData.Username
// 	password := requestData.Password

// 	if err := entities.UserDataVerification(username, password); err != nil {
// 		return UserResponse{}, errors.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
// 	}

// 	user, err := entities.CreateUser(username, password)
// 	if err != nil {
// 		return UserResponse{}, errCreateUser
// 	}

// 	responseWriter, ok := httputils.ContextWriter(ctx)
// 	if !ok {
// 		return UserResponse{}, errInternal
// 	}

// 	err = usecase.SetSession(responseWriter, user.ID)
// 	if err != nil {
// 		return UserResponse{}, errSetSession
// 	}
// 	return UserResponse{ID: user.ID, Username: user.Username}, nil
// }

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (UserResponse, error) {
	db, err := db.GetPostgres()
	if err != nil {
		logger.Logger().Error(err.Error())
	}

	username := requestData.Email
	password := requestData.Passwrd

	if err := entities.UserDataVerification(username, password); err != nil {
		return UserResponse{}, errors.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	user, err := entities.CreateUser(username, password)
	if err != nil {
		return UserResponse{}, errCreateUser
	}

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return UserResponse{}, errInternal
	}

	err = usecase.SetSession(responseWriter, user.ID)
	if err != nil {
		return UserResponse{}, errSetSession
	}

	dataStr := make(map[string]string)

	dataStr["email"] = user.Email
	dataStr["passwrd"] = user.Passwrd

	UserRepo := userRep.NewUserRepo(db)
	_, err = UserRepo.CreateUser(dataStr)

	if err != nil {
		return UserResponse{}, errCreateUser
	}

	return UserResponse{ID: user.ID, Username: user.Email}, nil
}
