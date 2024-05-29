package mocks

import (
	"context"
	"net/http"

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
