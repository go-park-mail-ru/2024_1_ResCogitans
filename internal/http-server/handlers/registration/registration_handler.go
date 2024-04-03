package registration

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
)

type RegistrationHandler struct {
	sessionUseCase usecase.AuthInterface
	userUseCase    usecase.UserInterface
}

func NewRegistrationHandler(sessionUseCase usecase.AuthInterface, userUseCase usecase.UserInterface) *RegistrationHandler {
	return &RegistrationHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (entities.UserResponse, httperrors.HttpError) {
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

	if h.userUseCase.IsUsernameTaken(username) {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, errors.New("Username is taken"))
	}

	if err := h.userUseCase.UserDataVerification(username, password); err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, err)
	}

	user, err := h.userUseCase.CreateUser(username, password)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}

	err = h.sessionUseCase.CreateSession(responseWriter, user.ID)
	if err != nil {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusInternalServerError, err)
	}
	return entities.UserResponse{Username: user.Username}, httperrors.HttpError{}
}
