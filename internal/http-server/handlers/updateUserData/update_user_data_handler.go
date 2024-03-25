package updateUserData

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
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

func (h *UpdateDataHandler) Update(ctx context.Context, requestData entities.User) (entities.UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	request, ok := httputils.GetRequestFromCtx(ctx)
	if !ok {
		return entities.UserResponse{}, fmt.Errorf("can't get request from context")
	}

	userID, err := h.useCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, err
	}

	user, err := entities.ChangeData(userID, username, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{ID: user.ID, Username: user.Username}, nil
}
