package authorization

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
)

type AuthorizationHandler struct {
	sessionUseCase usecase.AuthInterface
	userUseCase    usecase.UserInterface
}

func NewAuthorizationHandler(sessionUseCase usecase.AuthInterface, userUseCase usecase.UserInterface) *AuthorizationHandler {
	return &AuthorizationHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (entities.UserResponse, httperrors.HttpError) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	sessionID, _ := h.sessionUseCase.GetSession(request)
	if sessionID != 0 {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, errors.New("User is already authorized"))
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	if err := h.userUseCase.UserDataVerification(username, password); err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, err)
	}

	err = h.userUseCase.UserExists(username, password)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, err)
	}

	user, err := h.userUseCase.GetUserByUsername(username)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	err = h.sessionUseCase.CreateSession(responseWriter, user.ID)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	userResponse := entities.UserResponse{
		Username: user.Username,
	}

	return userResponse, httperrors.HttpError{}
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, _ entities.User) (entities.UserResponse, httperrors.HttpError) {
	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	_, err = h.sessionUseCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusForbidden, err)
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	err = h.sessionUseCase.ClearSession(responseWriter, request)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	return entities.UserResponse{}, httperrors.HttpError{}
}
