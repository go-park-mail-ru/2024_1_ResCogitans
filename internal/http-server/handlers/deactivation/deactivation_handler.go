package deactivation

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type DeactivationHandler struct {
	useCase usecase.AuthInterface
}

func NewDeactivationHandler(useCase usecase.AuthInterface) *DeactivationHandler {
	return &DeactivationHandler{
		useCase: useCase,
	}
}

func (h *DeactivationHandler) Deactivate(ctx context.Context, _ entities.User) (entities.UserResponse, httperrors.HttpError) {
	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	session, err := h.useCase.GetSession(request)
	if err != nil {
		errForbidden := httperrors.ErrForbidden
		errForbidden.Message = err
		return entities.UserResponse{}, errForbidden
	}

	err = entities.DeleteUser(session)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
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
