package registration

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
)

type RegistrationHandler struct{}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	if _, err := entities.UserDataVerification(username, password); err != nil {
		return UserResponse{}, errors.Wrap(err, "User data verification failed")
	}

	user, err := entities.CreateUser(username, password)
	if err != nil {
		return UserResponse{}, errors.Wrap(err, "Failed creating new profile")
	}

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return UserResponse{}, errors.New("Internal Error")
	}

	entities.SetSession(responseWriter, user.ID)
	return UserResponse{ID: user.ID, Username: user.Username}, nil
}
