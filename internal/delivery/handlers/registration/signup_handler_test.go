package registration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/mocks"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationHandler_SignUp(t *testing.T) {
	mockSessionUseCase := new(mocks.MockSessionUseCase)
	mockUserUseCase := new(mocks.MockUserUseCase)

	handler := registration.NewRegistrationHandler(mockSessionUseCase, mockUserUseCase)

	user := entities.User{Username: "san@boy.ru", Password: "ABC123abc123!"}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		assert.NoError(t, err)
	}

	// Создание запроса для добавления в контекст
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(jsonUser))
	if err != nil {
		assert.NoError(t, err)
	}

	// Создание ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	t.Run("Request without request key", func(t *testing.T) {
		userResponse, err := handler.SignUp(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		assert.Equal(t, "failed getting request", err.Error())
	})

	// Добавление запроса в контекст
	ctx := context.WithValue(req.Context(), httputils.HttpRequestKey, req)

	t.Run("Request without response writer key", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("CreateUser", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()

		userResponse, err := handler.SignUp(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, userResponse)
		assert.Equal(t, "failed getting response writer", err.Error())
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	// добавление ключа ответа в контекст
	ctx = context.WithValue(ctx, httputils.ResponseWriterKey, rr)

	t.Run("Successful registration", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("CreateUser", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(nil).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boy.ru"}, userResponse)
		assert.Equal(t, http.StatusOK, rr.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Cookie not found", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, httperrors.HttpError{Code: http.StatusBadRequest, Message: "cookie not found"}).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("CreateUser", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(nil).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boy.ru"}, userResponse)
		assert.Equal(t, http.StatusOK, rr.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error decoding cookie", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, httperrors.HttpError{Code: http.StatusInternalServerError, Message: "error decoding cookie"}).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("CreateUser", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(nil).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boy.ru"}, userResponse)
		assert.Equal(t, http.StatusOK, rr.Code)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error getting session", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, errors.New("unexpected error")).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error", err.Error())
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Authorized user", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(1, nil).Once()

		userResponse, err := handler.SignUp(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, "user is already authorized", err.Error())
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Email is taken", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(true, nil).Once()

		userResponse, err := handler.SignUp(ctx, user)

		assert.Error(t, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "username is taken", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Wrong username", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(httperrors.HttpError{Code: http.StatusBadRequest, Message: "email doesn't meet requirements"}).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.Error(t, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "email doesn't meet requirements", httpError.Message)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error creating user", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("CreateUser", ctx, "san@boy.ru", "ABC123abc123!").Return(errors.New("unexpected error")).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error", err.Error())
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error getting user by email", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("CreateUser", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{}, errors.New("unexpected error")).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error", err.Error())
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})

	mockUserUseCase.Mock.ExpectedCalls = nil
	mockSessionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error creating session", func(t *testing.T) {
		mockSessionUseCase.On("GetSession", ctx, req).Return(0, nil).Once()
		mockUserUseCase.On("IsEmailTaken", ctx, "san@boy.ru").Return(false, nil).Once()
		mockUserUseCase.On("UserDataVerification", "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("CreateUser", ctx, "san@boy.ru", "ABC123abc123!").Return(nil).Once()
		mockUserUseCase.On("GetUserByEmail", ctx, "san@boy.ru").Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil).Once()
		mockSessionUseCase.On("CreateSession", ctx, rr, 1).Return(errors.New("unexpected error")).Once()

		userResponse, err := handler.SignUp(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error", err.Error())
		assert.Equal(t, entities.UserResponse{}, userResponse)

		mockSessionUseCase.AssertExpectations(t)
		mockUserUseCase.AssertExpectations(t)
	})
}
