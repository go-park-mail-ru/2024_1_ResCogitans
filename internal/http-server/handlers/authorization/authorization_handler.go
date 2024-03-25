package authorization

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type AuthorizationHandler struct {
	useCase usecase.AuthInterface
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
	errInternal = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "internal Error",
	}
	errSessionNotSet = errors.HttpError{
		Code:    http.StatusUnauthorized,
		Message: "session is not set",
	}
)

func NewAuthorizationHandler(useCase usecase.AuthInterface) *AuthorizationHandler {
	return &AuthorizationHandler{
		useCase: useCase,
	}
}

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (entities.UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	responseWriter, ok := httputils.GetResponseWriterFromCtx(ctx)
	if !ok {
		return entities.UserResponse{}, errInternal
	}

	if err := entities.UserDataVerification(username, password); err != nil {
		return entities.UserResponse{}, errors.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	user, err := entities.GetUserByUsername(username)
	if err != nil {
		return entities.UserResponse{}, errLoginUser
	}

	if ok = entities.IsAuthenticated(username, password); !ok {
		return entities.UserResponse{}, errLoginUser
	}

	err = h.useCase.SetSession(responseWriter, user.ID)

	if err != nil {
		return entities.UserResponse{}, errSetSession
	}

	userResponse := entities.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, requestData entities.User) (entities.UserResponse, error) {
	request, ok := httputils.GetRequestFromCtx(ctx)
	if !ok {
		return entities.UserResponse{}, errInternal
	}

	responseWriter, ok := httputils.GetResponseWriterFromCtx(ctx)
	if !ok {
		return entities.UserResponse{}, errInternal
	}

	userID, err := h.useCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, errSetSession
	}

	if userID == 0 {
		return entities.UserResponse{}, errSessionNotSet
	}

	err = h.useCase.ClearSession(responseWriter, request)
	if err != nil {
		return entities.UserResponse{}, errClearSession
	}

	return entities.UserResponse{}, nil
}
