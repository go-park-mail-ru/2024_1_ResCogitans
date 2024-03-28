package registration

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
)

type RegistrationHandler struct {
	useCase usecase.AuthInterface
}

func NewRegistrationHandler(useCase usecase.AuthInterface) *RegistrationHandler {
	return &RegistrationHandler{
		useCase: useCase,
	}
}

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (entities.UserResponse, httperrors.HttpError) {
	username := requestData.Username
	password := requestData.Password

	if _, err := entities.GetUserByUsername(username); err == nil {
		errBadRequest := httperrors.ErrBadRequest
		errBadRequest.Message = errors.New("User with this username is already exists")
		return entities.UserResponse{}, errBadRequest
	}

	if err := entities.UserDataVerification(username, password); err != nil {
		errBadRequest := httperrors.ErrBadRequest
		errBadRequest.Message = err
		return entities.UserResponse{}, errBadRequest
	}

	user, err := entities.CreateUser(username, password)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}

	err = h.useCase.SetSession(responseWriter, user.ID)
	if err != nil {
		errInternal := httperrors.ErrInternal
		errInternal.Message = err
		return entities.UserResponse{}, errInternal
	}
	return entities.UserResponse{ID: user.ID, Username: user.Username}, httperrors.HttpError{}
}
