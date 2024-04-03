package deactivation

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type DeactivationHandler struct {
	sessionUseCase usecase.AuthInterface
	userUseCase    usecase.UserInterface
}

func NewDeactivationHandler(sessionUseCase usecase.AuthInterface, userUseCase usecase.UserInterface) *DeactivationHandler {
	return &DeactivationHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

func (h *DeactivationHandler) Deactivate(ctx context.Context, _ entities.User) (entities.UserResponse, httperrors.HttpError) {
	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	session, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusForbidden, err)
	}

	err = h.userUseCase.DeleteUser(session)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
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
