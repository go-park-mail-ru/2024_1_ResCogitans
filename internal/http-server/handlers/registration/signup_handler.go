package registration

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type Registration struct{}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

func (h *Registration) SignUp(ctx context.Context, _ entities.User) (Response, error) {
	requestData, ok := ctx.Value("requestData").(entities.User)
	logger.Logger().DebugContext(ctx, "str")
	if !ok {
		return Response{Status: http.StatusBadRequest, Message: "requestData not found in context"}, nil
	}

	username := requestData.Username
	password := requestData.Password

	if status, err := entities.UserDataVerification(username, password); err != nil {
		return Response{Status: status, Message: err.Error()}, err
	}

	_, err := entities.CreateUser(username, password)
	if err != nil {
		return Response{Status: http.StatusBadRequest, Message: "Failed creating new profile"}, err
	}

	return Response{Status: http.StatusCreated}, nil
}
