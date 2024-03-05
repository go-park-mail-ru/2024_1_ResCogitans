package registration

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/pkg/errors"
)

type Registration struct{}

func (h *Registration) SignUp(ctx context.Context, _ entities.User) (*entities.User, error) {
	requestData, ok := ctx.Value("requestData").(entities.User)
	logger.Logger().DebugContext(ctx, "str")
	if !ok {
		return nil, errors.New("requestData not found in context")
	}

	username := requestData.Username
	password := requestData.Password

	if err := entities.UserDataVerification(username, password); err != nil {
		return nil, err
	}

	newUser, err := entities.CreateUser(username, password)
	if err != nil {
		return nil, errors.New("Failed creating new profile")
	}

	return &newUser, nil
}
