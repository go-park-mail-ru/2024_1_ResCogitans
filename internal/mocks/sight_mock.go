package mocks

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MockSightUseCase struct {
	mock.Mock
}

func (m *MockSightUseCase) GetSightByID(ctx context.Context, sightID int) (entities.Sight, error) {
	args := m.Called(ctx, sightID)
	return args.Get(0).(entities.Sight), args.Error(1)
}

func (m *MockSightUseCase) GetCommentsBySightID(ctx context.Context, commentID int) ([]entities.Comment, error) {
	args := m.Called(ctx, commentID)
	return args.Get(0).([]entities.Comment), args.Error(1)
}

func (m *MockSightUseCase) GetCommentsByUserID(ctx context.Context, userID int) ([]entities.Comment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entities.Comment), args.Error(1)
}

func (m *MockSightUseCase) GetSightsList(ctx context.Context) ([]entities.Sight, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Sight), args.Error(1)
}

func (m *MockSightUseCase) SearchSights(ctx context.Context, str string) (entities.Sights, error) {
	args := m.Called(ctx, str)
	return args.Get(0).(entities.Sights), args.Error(1)
}
