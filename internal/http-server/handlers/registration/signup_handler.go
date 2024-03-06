package registration

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RegistrationHandler struct{}

type RegResponse struct {
	User      entities.User
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

	responseWriter, ok := authorization.ContextWriter(ctx)
	if !ok {
		return RegResponse{}, errors.New("Internal Error")
	}

	sessionID := uuid.New().String()
	authorization.SessionStore[sessionID] = username
	authorization.SetSession(sessionID, responseWriter)

	return RegResponse{User: user, SessionID: sessionID}, nil
}
