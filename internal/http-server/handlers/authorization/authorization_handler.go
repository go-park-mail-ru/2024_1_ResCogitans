package authorization

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type AuthorizationHandler struct{}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

var (
	errLoginUser = &errors.HttpError{
		Code:    http.StatusBadRequest,
		Message: "failed authorize",
	}
	errSetSession = &errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed setting session",
	}
	errClearSession = &errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed clearing session",
	}
	errInternal = &errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "internal Error",
	}
	errSessionNotSet = &errors.HttpError{
		Code:    http.StatusUnauthorized,
		Message: "session is not set",
	}
)

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (UserResponse, *errors.HttpError) {
	username := requestData.Username
	password := requestData.Password

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return UserResponse{}, errInternal
	}

	user, err := entities.GetUserByUsername(username)
	if err != nil {
		return UserResponse{}, errLoginUser
	}

	if ok = entities.IsAuthenticated(username, password); !ok {
		return UserResponse{}, errLoginUser
	}

	err = usecase.SetSession(responseWriter, user.ID)
	if err != nil {
		return UserResponse{}, errSetSession
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, requestData entities.User) (UserResponse, *errors.HttpError) {
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
