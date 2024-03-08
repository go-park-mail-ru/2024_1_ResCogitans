package authorization

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
	"net/http"
)

type AuthorizationHandler struct{}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return UserResponse{}, errors.New("Internal Error")
	}

	if entities.UserValidation(username, password) {
		user, err := entities.GetUserByUsername(username)
		if err != nil {
			return UserResponse{}, errors.Wrap(err, "Problem with searching for a profile by username")
		}

		entities.SetSession(responseWriter, user.ID)

		userResponse := UserResponse{
			ID:       user.ID,
			Username: user.Username,
		}

		return userResponse, nil
	}

	return UserResponse{}, errors.New("Authorization failed")
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, requestData entities.User) (UserResponse, error) {
	request, ok := ctx.Value(httputils.HttpRequestKey).(*http.Request)
	if !ok {
		return UserResponse{}, errors.New("failed getting http.request")
	}

	responseWriter, ok := httputils.ContextWriter(ctx)
	if !ok {
		return UserResponse{}, errors.New("Internal Error")
	}

	userID := entities.GetSession(request)

	if userID == 0 {
		return UserResponse{}, errors.New("session is not set")
	}

	entities.ClearSession(responseWriter, request)

	return UserResponse{}, nil
}
