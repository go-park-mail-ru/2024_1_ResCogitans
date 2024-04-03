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
	sessionUseCase usecase.AuthInterface
	userUseCase    usecase.UserInterface
}

func NewUpdateDataHandler(sessionUseCase usecase.AuthInterface, userUseCase usecase.UserInterface) *UpdateDataHandler {
	return &UpdateDataHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

func (h *UpdateDataHandler) Update(ctx context.Context, requestData entities.User) (entities.UserResponse, httperrors.HttpError) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	userID, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusForbidden, err)
	}

	err = h.userUseCase.UserDataVerification(username, password)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, err)
	}

	user, err := h.userUseCase.ChangeData(userID, username, password)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	return entities.UserResponse{Username: user.Username}, httperrors.HttpError{}
}
