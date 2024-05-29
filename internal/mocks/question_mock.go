package mocks

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/stretchr/testify/mock"
)

type QuestionMockUseCase struct {
	mock.Mock
}

func (m *QuestionMockUseCase) CreateReview(ctx context.Context, userID int, review entities.Review) error {
	args := m.Called(ctx, userID, review)
	return args.Error(0)
}

func (m *QuestionMockUseCase) GetQuestions(ctx context.Context) ([]entities.QuestionResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.QuestionResponse), args.Error(1)
}

func (m *QuestionMockUseCase) CheckReview(ctx context.Context, userID int) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *QuestionMockUseCase) SetStat(ctx context.Context, userID int) ([]entities.Statistic, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entities.Statistic), args.Error(1)
}
