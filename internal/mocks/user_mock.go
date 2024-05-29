package mocks

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/stretchr/testify/mock"
)

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
