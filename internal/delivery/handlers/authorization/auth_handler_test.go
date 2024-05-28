package authorization_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/mocks"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Тесты для авторизации
func TestAuthorizationHandler_Authorize(t *testing.T) {
	mockSessionUseCase := new(mocks.MockSessionUseCase)
	mockUserUseCase := new(mocks.MockUserUseCase)

	handler := authorization.NewAuthorizationHandler(mockSessionUseCase, mockUserUseCase)

	user := entities.User{Username: "san@boy.ru", Password: "ABC123abc123!"}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		assert.NoError(t, err)
	}

	// Создание запроса для добавления в контекст
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonUser))
	if err != nil {
		assert.NoError(t, err)
	}

	// Создание ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	t.Run("Request without request key", func(t *testing.T) {
		userResponse, err := handler.Authorize(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		assert.Equal(t, "failed getting request", err.Error())
	})

	// Добавление запроса в контекст
	ctx := context.WithValue(req.Context(), httputils.HttpRequestKey, req)

	t.Run("Request without response writer key", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("UserExists", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		assert.Equal(t, "failed getting response writer", err.Error())
	})

	// добавление ключа ответа в контекст
	ctx = context.WithValue(ctx, httputils.ResponseWriterKey, rr)

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Successful authorization", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("UserExists", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(nil).Once()

		userResponse, err := handler.Authorize(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boy.ru"}, userResponse)
		assert.Equal(t, http.StatusOK, rr.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Wrong username", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("UserExists", ctx, "san@boy.ru", "ABC123abc123!").Return(httperrors.HttpError{Code: http.StatusBadRequest, Message: "user not found"}).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		assert.IsType(t, httperrors.HttpError{}, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "user not found", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Incorrect username", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(httperrors.HttpError{Code: http.StatusBadRequest, Message: "email doesn't meet requirements"}).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		assert.IsType(t, httperrors.HttpError{}, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "email doesn't meet requirements", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Incorrect password", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(httperrors.HttpError{Code: http.StatusBadRequest, Message: "password doesn't meet requirements"}).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		assert.IsType(t, httperrors.HttpError{}, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "password doesn't meet requirements", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Failed getting user by username", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("UserExists", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{}, httperrors.HttpError{Code: http.StatusBadRequest, Message: "user not found"}).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		assert.IsType(t, httperrors.HttpError{}, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "user not found", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Failed creating session", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("UserExists", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(errors.New("unexpected error")).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "unexpected error", httpError.Message)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Cookie not found", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, httperrors.HttpError{Code: http.StatusBadRequest, Message: "cookie not found"}).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("UserExists", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(nil).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boy.ru"}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error decoding cookie", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, errors.New("error decoding cookie")).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("UserExists", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(nil).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boy.ru"}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Unexpected error getting session", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, errors.New("unexpected error")).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "unexpected error", httpError.Message)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("User is already authorised", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", mock.Anything, mock.Anything).Return(1, nil).Once()

		userResponse, err := handler.Authorize(ctx, user)

		assert.Error(t, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "user is already authorized", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})
}

// Тесты для выхода из аккаунта
func TestAuthorizationHandler_LogOut(t *testing.T) {
	mockSessionUseCase := new(mocks.MockSessionUseCase)
	mockUserUseCase := new(mocks.MockUserUseCase)

	handler := authorization.NewAuthorizationHandler(mockSessionUseCase, mockUserUseCase)

	user := entities.User{}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		assert.NoError(t, err)
	}

	// Создание запроса
	req, err := http.NewRequest("POST", "/api/logout", bytes.NewBuffer(jsonUser))
	if err != nil {
		assert.NoError(t, err)
	}

	// Создание ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	t.Run("Request without request key", func(t *testing.T) {
		userResponse, err := handler.LogOut(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		assert.Equal(t, "failed getting request", err.Error())
	})

	// Добавление запроса в контекст
	ctx := context.WithValue(req.Context(), httputils.HttpRequestKey, req)

	t.Run("Request without response writer key", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()

		userResponse, err := handler.LogOut(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		assert.Equal(t, "failed getting response writer", err.Error())

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	// добавление ключа ответа в контекст
	ctx = context.WithValue(ctx, httputils.ResponseWriterKey, rr)

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Cookie not found", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, httperrors.HttpError{Code: http.StatusUnauthorized, Message: "cookie not found"}).Once()

		userResponse, err := handler.LogOut(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "cookie not found", httpError.Message)
		assert.Equal(t, http.StatusUnauthorized, httpError.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error clearing session", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(1, nil).Once()
		mockSessionUseCase.On("ClearSession", ctx, rr, req).Return(errors.New("error clearing session")).Once()

		userResponse, err := handler.LogOut(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "error clearing session", httpError.Message)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Success log out", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(1, nil).Once()
		mockSessionUseCase.On("ClearSession", ctx, rr, req).Return(nil).Once()

		userResponse, err := handler.LogOut(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})
}

// Тесты для обновления пароля
func TestAuthorizationHandler_UpdatePassword(t *testing.T) {
	mockSessionUseCase := new(mocks.MockSessionUseCase)
	mockUserUseCase := new(mocks.MockUserUseCase)

	handler := authorization.NewAuthorizationHandler(mockSessionUseCase, mockUserUseCase)

	user := entities.User{Username: "san@boys.ru", Password: "ABC123abc1234!"}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		assert.NoError(t, err)
	}

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonUser))
	if err != nil {
		assert.NoError(t, err)
	}

	t.Run("Request without request key", func(t *testing.T) {
		userResponse, err := handler.UpdatePassword(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		assert.Equal(t, "failed getting request", err.Error())

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	// Добавление запроса в контекст
	ctx := context.WithValue(req.Context(), httputils.HttpRequestKey, req)

	t.Run("Cookie nit found", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, httperrors.HttpError{Code: http.StatusBadRequest, Message: "cookie not found"}).Once()

		userResponse, err := handler.UpdatePassword(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "cookie not found", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Incorrect username", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(1, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boys.ru", "ABC123abc1234!").Return(httperrors.HttpError{Code: http.StatusBadRequest, Message: "email doesn't meet requirements"}).Once()

		userResponse, err := handler.UpdatePassword(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "email doesn't meet requirements", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error changing password", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(1, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boys.ru", "ABC123abc1234!").Return(nil).Once()
		mockUserUseCase.On("ChangePassword", ctx, 1, "ABC123abc1234!").Return(entities.User{}, errors.New("failed changing password")).Once()

		userResponse, err := handler.UpdatePassword(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "failed changing password", httpError.Message)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Success changing password", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(1, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boys.ru", "ABC123abc1234!").Return(nil).Once()
		mockUserUseCase.On("ChangePassword", ctx, 1, "ABC123abc1234!").Return(entities.User{ID: 1, Username: "san@boys.ru"}, nil).Once()

		userResponse, err := handler.UpdatePassword(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boys.ru"}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})
}
