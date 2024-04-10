package delivery

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

type AuthorizationHandler struct{}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

var (
	errLoginUser = errors.HttpError{
		Code:    http.StatusBadRequest,
		Message: "failed authorize",
	}
	errSetSession = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed setting session",
	}
	errClearSession = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed clearing session",
	}
	errSessionNotSet = errors.HttpError{
		Code:    http.StatusUnauthorized,
		Message: "session is not set",
	}
)

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (UserResponse, error) {
	username := requestData.Email
	password := requestData.Passwrd

	db, err := db.GetPostgres()
	if err != nil {
		logger.Logger().Error(err.Error())
	}

	dataStr := make(map[string]string)

	dataStr["email"] = username
	dataStr["passwrd"] = password

	UserRepo := userRep.NewUserRepo(db)
	user, err := UserRepo.AuthorizeUser(dataStr)
	if err != nil {
		return UserResponse{}, errLoginUser
	}

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return UserResponse{}, errInternal
	}

	err = usecase.SetSession(responseWriter, user.ID)
	if err != nil {
		return UserResponse{}, errSetSession
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Email,
	}

	return userResponse, nil
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, requestData entities.User) (UserResponse, error) {
	request, ok := httputils.HttpRequest(ctx)
	if !ok {
		return UserResponse{}, errInternal
	}

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return UserResponse{}, errInternal
	}

	userID := usecase.GetSession(request)

	if userID == 0 {
		return UserResponse{}, errSessionNotSet
	}

	err := usecase.ClearSession(responseWriter, request)
	if err != nil {
		return UserResponse{}, errClearSession
	}

	return UserResponse{}, nil
}
