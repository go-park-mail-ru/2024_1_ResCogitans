package updateUserData

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type UpdateDataHandler struct {
	useCase usecase.AuthInterface
}

func NewUpdateDataHandler(useCase usecase.AuthInterface) *UpdateDataHandler {
	return &UpdateDataHandler{
		useCase: useCase,
	}
}

func (h *UpdateDataHandler) Update(ctx context.Context, requestData entities.User) (entities.UserResponse, httperrors.HttpError) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	userID, err := h.useCase.GetSession(request)
	if err != nil {
		errUnauthorized := httperrors.ErrUnauthorized
		errUnauthorized.Message = err
		return entities.UserResponse{}, errUnauthorized
	}

	user, err := entities.ChangeData(userID, username, password)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	return entities.UserResponse{ID: user.ID, Username: user.Username}, httperrors.HttpError{}
}
