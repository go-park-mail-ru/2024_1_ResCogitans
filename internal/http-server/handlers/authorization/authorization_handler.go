package authorization

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
)

type AuthorizationHandler struct {
	useCase usecase.AuthInterface
}

func NewAuthorizationHandler(useCase usecase.AuthInterface) *AuthorizationHandler {
	return &AuthorizationHandler{
		useCase: useCase,
	}
}

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (entities.UserResponse, httperrors.HttpError) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	session, _ := h.useCase.GetSession(request)
	if session != 0 {
		errForbidden := httperrors.ErrForbidden
		errForbidden.Message = errors.New("User is already authorized")
		return entities.UserResponse{}, errForbidden
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		errBadRequest := httperrors.ErrBadRequest
		errBadRequest.Message = err
		return entities.UserResponse{}, errBadRequest
	}

	if err := entities.UserDataVerification(username, password); err != nil {
		errBadRequest := httperrors.ErrBadRequest
		errBadRequest.Message = err
		return entities.UserResponse{}, errBadRequest
	}

	if err := entities.UserExists(username, password); err != nil {
		errBadRequest := httperrors.ErrBadRequest
		errBadRequest.Message = err
		return entities.UserResponse{}, errBadRequest
	}

	user, err := entities.GetUserByUsername(username)
	if err != nil {
		errBadRequest := httperrors.ErrBadRequest
		errBadRequest.Message = err
		return entities.UserResponse{}, errBadRequest
	}

	err = h.useCase.SetSession(responseWriter, user.ID)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	userResponse := entities.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, httperrors.HttpError{}
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, _ entities.User) (entities.UserResponse, httperrors.HttpError) {
	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	userID, err := h.useCase.GetSession(request)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	if userID == 0 {
		errInternal := httperrors.ErrInternal
		errInternal.Message = errors.New("Session not found")
		return entities.UserResponse{}, errInternal
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	err = h.useCase.ClearSession(responseWriter, request)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	return entities.UserResponse{}, httperrors.HttpError{}
}
