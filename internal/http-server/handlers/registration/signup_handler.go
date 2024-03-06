package registration

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RegistrationHandler struct{}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type RegResponse struct {
	User      UserResponse
	SessionID string
}

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (RegResponse, error) {
	username := requestData.Username
	password := requestData.Password

	if _, err := entities.UserDataVerification(username, password); err != nil {
		return RegResponse{}, errors.Wrap(err, "User data verification failed")
	}

	user, err := entities.CreateUser(username, password)
	if err != nil {
		return RegResponse{}, errors.Wrap(err, "Failed creating new profile")
	}

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return RegResponse{}, errors.New("Internal Error")
	}

	sessionID := uuid.New().String()
	authorization.SessionStore[sessionID] = username

	if err := authorization.SetSession(sessionID, responseWriter); err != nil {
		return RegResponse{}, errors.Wrap(err, "failed to set session cookie")
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	return RegResponse{User: userResponse, SessionID: sessionID}, nil
}
