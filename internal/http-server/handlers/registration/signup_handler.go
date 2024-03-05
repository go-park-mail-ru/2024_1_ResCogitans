package registration

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/pkg/errors"
)

type Registration struct{}

func (h *Registration) SignUp(ctx context.Context, _ entities.User) (*entities.User, error) {
	fmt.Println("ABOBIUM")
	username := ctx.Value("username").(string)
	password := ctx.Value("password").(string)

	if err := entities.UserDataVerification(username, password); err != nil {
		return nil, err
	}

	newUser, err := entities.CreateUser(username, password)
	if err != nil {
		return nil, errors.New("Failed creating new profile")
	}

	return newUser, nil
}
