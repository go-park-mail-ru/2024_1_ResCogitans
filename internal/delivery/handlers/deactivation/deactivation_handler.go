package deactivation

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type DeactivationHandler struct {
	sessionUseCase usecase.SessionInterface
	userUseCase    usecase.UserUseCaseInterface
}

func NewDeactivationHandler(sessionUseCase usecase.SessionInterface, userUseCase usecase.UserUseCaseInterface) *DeactivationHandler {
	return &DeactivationHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

// Deactivate godoc
// @Summary Удаление аккаунта
// @Description Удаляет из базы данных user и profile
// @Tags Деактивация
// @Accept json
// @Produce json
// @Success 200 {object} entities.UserResponse
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id}/delete [get]
func (h *DeactivationHandler) Deactivate(ctx context.Context, _ entities.UserRequest) (entities.UserResponse, error) {
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
