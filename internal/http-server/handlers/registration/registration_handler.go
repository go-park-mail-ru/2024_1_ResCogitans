package registration

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type RegistrationHandler struct {
	sessionUseCase usecase.SessionInterface
	userUseCase    usecase.UserInterface
}

func NewRegistrationHandler(sessionUseCase usecase.SessionInterface, userUseCase usecase.UserInterface) *RegistrationHandler {
	return &RegistrationHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

func (h *RegistrationHandler) SignUp(ctx context.Context, requestData entities.User) (entities.UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	println("1")

	sessionID, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		httpError := httperrors.UnwrapHttpError(err)
		if httpError.Message != "Session not found" && httpError.Message != "Cookie not found" {
			return entities.UserResponse{}, httpError
		}
	}

	println("2")

	if sessionID != 0 {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, "User is already authorized")
	}

	println("3")

	if h.userUseCase.IsUsernameTaken(username) {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, "Username is taken")
	}

	println("4")

	if err := h.userUseCase.UserDataVerification(username, password); err != nil {
		return entities.UserResponse{}, err
	}

	println("5")

	user, err := h.userUseCase.CreateUser(username, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	println("6")

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	println("7")

	err = h.sessionUseCase.CreateSession(responseWriter, user.ID)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{Username: user.Username}, nil
}
