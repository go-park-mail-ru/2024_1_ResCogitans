package registration

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type RegistrationHandler struct {
	useCase usecase.AuthInterface
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

var (
	errCreateUser = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed creating new profile",
	}
	errSetSession = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed setting session",
	}
	errInternal = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "internal Error",
	}
)

func NewRegistrationHandler(useCase usecase.AuthInterface) *RegistrationHandler {
	return &RegistrationHandler{
		useCase: useCase,
	}
}

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	if err := entities.UserDataVerification(username, password); err != nil {
		return UserResponse{}, errors.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err := entities.UserExists(username); err != nil {
		return UserResponse{}, errors.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	user, err := entities.CreateUser(username, password)
	if err != nil {
		return UserResponse{}, errCreateUser
	}

	responseWriter, ok := httputils.GetResponseWriterFromCtx(ctx)
	if !ok {
		return UserResponse{}, errInternal
	}

	err = h.useCase.SetSession(responseWriter, user.ID)
	if err != nil {
		return UserResponse{}, errSetSession
	}
	return UserResponse{ID: user.ID, Username: user.Username}, nil
}
