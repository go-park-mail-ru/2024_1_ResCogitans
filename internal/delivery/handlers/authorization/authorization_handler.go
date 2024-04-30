package authorization

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type AuthorizationHandler struct {
	sessionUseCase usecase.SessionInterface
	userUseCase    usecase.UserUseCaseInterface
}

func NewAuthorizationHandler(sessionUseCase usecase.SessionInterface, userUseCase usecase.UserUseCaseInterface) *AuthorizationHandler {
	return &AuthorizationHandler{
		sessionUseCase: sessionUseCase,
		userUseCase:    userUseCase,
	}
}

// Authorize godoc
// @Summary Авторизация пользователя
// @Description Авторизует пользователя и возвращает информацию о нем
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param user body entities.UserRequest true "Информация о пользователе, которая нужна при авторизации"
// @Success 200 {object} entities.UserResponse
// @Failure 400 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/login [post]
func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.UserRequest) (entities.UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	sessionID, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		if !httperrors.IsHttpError(err) {
			return entities.UserResponse{}, err
		}

		httpError := httperrors.UnwrapHttpError(err)
		if httpError.Message != "Cookie not found" && httpError.Message != "Error decoding cookie" {
			return entities.UserResponse{}, httpError
		}
	}

	if sessionID != 0 {
		return entities.UserResponse{}, httperrors.NewHttpError(http.StatusBadRequest, "User is already authorized")
	}

	if err := h.userUseCase.UserDataVerification(username, password); err != nil {
		return entities.UserResponse{}, err
	}

	err = h.userUseCase.UserExists(username, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	user, err := h.userUseCase.GetUserByEmail(username)
	if err != nil {
		return entities.UserResponse{}, err
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	err = h.sessionUseCase.CreateSession(responseWriter, user.ID)
	if err != nil {
		return entities.UserResponse{}, err
	}

	userResponse := entities.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}

// LogOut godoc
// @Summary Выход из аккаунта
// @Description Выход авторизованного пользователя из аккаунта
// @Tags Авторизация
// @Accept json
// @Produce json
// @Success 200 {object} entities.UserResponse
// @Failure 400 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/logout [post]
func (h *AuthorizationHandler) LogOut(ctx context.Context, _ entities.UserRequest) (entities.UserResponse, error) {
	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	_, err = h.sessionUseCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, err
	}

	responseWriter, err := httputils.GetResponseWriterFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	err = h.sessionUseCase.ClearSession(responseWriter, request)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{}, nil
}

// UpdatePassword godoc
// @Summary Смена пароля
// @Description Смена пароля авторизованного пользователя. Необходим только новый пароль
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param user body entities.UserRequest true "Ник пользователя и новый пароль"
// @Success 200 {object} entities.UserResponse
// @Failure 400 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id}/reset_password [post]
func (h *AuthorizationHandler) UpdatePassword(ctx context.Context, requestData entities.UserRequest) (entities.UserResponse, error) {
	username := requestData.Username
	password := requestData.Password

	request, err := httputils.GetRequestFromCtx(ctx)
	if err != nil {
		return entities.UserResponse{}, err
	}

	userID, err := h.sessionUseCase.GetSession(request)
	if err != nil {
		return entities.UserResponse{}, err
	}

	err = h.userUseCase.UserDataVerification(username, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	user, err := h.userUseCase.ChangePassword(userID, password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{Username: user.Username, ID: user.ID}, nil
}
