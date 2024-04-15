package updateUserData

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type UpdateDataHandler struct {
	sessionUseCase usecase.SessionInterface
	userUseCase    usecase.UserInterface
}

func NewUpdateDataHandler(sessionUseCase usecase.SessionInterface, userUseCase usecase.UserInterface) *UpdateDataHandler {
	return &UpdateDataHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

func (h *UpdateDataHandler) Update(ctx context.Context, requestData entities.User) (entities.UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	userID, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, err
	}

	if h.userUseCase.IsUsernameTaken(username) {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, "Username is taken")
	}

	err = h.userUseCase.UserDataVerification(username, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	user, err := h.userUseCase.ChangeData(userID, username, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{Username: user.Username}, nil
}
