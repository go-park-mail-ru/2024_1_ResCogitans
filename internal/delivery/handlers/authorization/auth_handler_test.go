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
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSessionUseCase struct {
	mock.Mock
}

func (m *MockSessionUseCase) CreateSession(ctx context.Context, w http.ResponseWriter, userID int) error {
	args := m.Called(ctx, w, userID)
	return args.Error(0)
}

func (m *MockSessionUseCase) GetSession(ctx context.Context, r *http.Request) (int, error) {
	args := m.Called(ctx, r)
	return args.Int(0), args.Error(1)
}

func (m *MockSessionUseCase) ClearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	args := m.Called(ctx, w, r)
	return args.Error(0)
}

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) CreateUser(ctx context.Context, email string, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByID(ctx context.Context, userID int) (entities.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserUseCase) UserDataVerification(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *MockUserUseCase) ChangePassword(ctx context.Context, userID int, password string) (entities.User, error) {
	args := m.Called(ctx, userID, password)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserUseCase) UserExists(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

func (m *MockUserUseCase) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func TestAuthorizationHandler_Authorize(t *testing.T) {
	println("________________________________________________________________________")
	mockSessionUseCase := new(MockSessionUseCase)
	mockUserUseCase := new(MockUserUseCase)

	handler := authorization.NewAuthorizationHandler(mockSessionUseCase, mockUserUseCase)

	// Тест успешной авторизации
	mockSessionUseCase.On("GetSession", mock.Anything, mock.Anything).Return(0, nil)
	mockUserUseCase.On("UserDataVerification", mock.Anything, mock.Anything).Return(nil)
	mockUserUseCase.On("UserExists", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockUserUseCase.On("GetUserByEmail", mock.Anything, mock.Anything).Return(entities.User{ID: 1, Username: "san@boy.ru"}, nil)
	mockSessionUseCase.On("CreateSession", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Создаем тестовый запрос
	user := entities.User{Username: "san@boy.ru", Password: "ABC123abc123!"}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Добавляем запрос в контекст
	ctx := context.WithValue(req.Context(), httputils.HttpRequestKey, req)
	ctx = context.WithValue(ctx, httputils.ResponseWriterKey, rr)

	userResponse, err := handler.Authorize(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, entities.UserResponse{ID: 1, Username: "san@boy.ru"}, userResponse)
	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что методы были вызваны с правильными аргументами
	mockSessionUseCase.AssertExpectations(t)
	mockUserUseCase.AssertExpectations(t)
}

// Тесты для LogOut
//func TestAuthorizationHandler_LogOut(t *testing.T) {
//	mockSessionUseCase := new(MockSessionUseCase)
//	mockUserUseCase := new(MockUserUseCase)
//
//	handler := authorization.NewAuthorizationHandler(mockSessionUseCase, mockUserUseCase)
//
//	// Тест успешного выхода из системы
//	mockSessionUseCase.On("GetSession", mock.Anything, mock.Anything).Return(1, nil)
//	mockSessionUseCase.On("ClearSession", mock.Anything, mock.Anything, mock.Anything).Return(nil)
//
//	_, err := handler.LogOut(context.Background(), entities.User{})
//	assert.NoError(t, err)
//
//	// Добавьте дополнительные тесты для других случаев, таких как ошибки получения сессии и ошибки очистки сессии
//}
//
//// Тесты для UpdatePassword
//func TestAuthorizationHandler_UpdatePassword(t *testing.T) {
//	mockSessionUseCase := new(MockSessionUseCase)
//	mockUserUseCase := new(MockUserUseCase)
//
//	handler := authorization.NewAuthorizationHandler(mockSessionUseCase, mockUserUseCase)
//
//	// Тест успешного обновления пароля
//	mockSessionUseCase.On("GetSession", mock.Anything, mock.Anything).Return(1, nil)
//	mockUserUseCase.On("UserDataVerification", mock.Anything, mock.Anything).Return(nil)
//	mockUserUseCase.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(entities.User{ID: 1, Username: "testuser"}, nil)
//
//	userResponse, err := handler.UpdatePassword(context.Background(), entities.User{Username: "testuser", Password: "newpassword"})
//	assert.NoError(t, err)
//	assert.Equal(t, entities.UserResponse{ID: 1, Username: "testuser"}, userResponse)
//
//	// Добавьте дополнительные тесты для других случаев, таких как ошибки валидации данных пользователя, ошибки изменения пароля и т.д.
//}
