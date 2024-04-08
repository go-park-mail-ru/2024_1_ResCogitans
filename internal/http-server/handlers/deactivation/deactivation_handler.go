package deactivation

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
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

func (h *DeactivationHandler) Deactivate(ctx context.Context, _ entities.User) (entities.UserResponse, error) {
	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	session, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, err
	}

	err = h.userUseCase.DeleteUser(session)
	if err != nil {
		return entities.UserResponse{}, err
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	err = h.sessionUseCase.ClearSession(responseWriter, request)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{}, nil
}
