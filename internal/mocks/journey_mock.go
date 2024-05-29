package mocks

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/stretchr/testify/mock"
)

type JourneyMockUseCase struct {
	mock.Mock
}

func (m *JourneyMockUseCase) CreateJourney(ctx context.Context, journey entities.Journey) (entities.Journey, error) {
	args := m.Called(ctx, journey)
	return args.Get(0).(entities.Journey), args.Error(1)
}

func (m *JourneyMockUseCase) DeleteJourneyByID(ctx context.Context, journeyID int) error {
	args := m.Called(ctx, journeyID)
	return args.Error(0)
}

func (m *JourneyMockUseCase) GetJourneys(ctx context.Context, userID int) ([]entities.Journey, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entities.Journey), args.Error(1)
}

func (m *JourneyMockUseCase) AddJourneySight(ctx context.Context, journeyID int, ids []int) error {
	args := m.Called(ctx, journeyID, ids)
	return args.Error(0)
}

func (m *JourneyMockUseCase) EditJourney(ctx context.Context, journeyID int, name, description string) error {
	args := m.Called(ctx, journeyID, name, description)
	return args.Error(0)
}

func (m *JourneyMockUseCase) DeleteJourneySight(ctx context.Context, journeyID int, sight entities.JourneySight) error {
	args := m.Called(ctx, journeyID, sight)
	return args.Error(0)
}

func (m *JourneyMockUseCase) GetJourneySights(ctx context.Context, journeyID int) ([]entities.Sight, error) {
	args := m.Called(ctx, journeyID)
	return args.Get(0).([]entities.Sight), args.Error(1)
}

func (m *JourneyMockUseCase) GetJourney(ctx context.Context, journeyID int) (entities.Journey, error) {
	args := m.Called(ctx, journeyID)
	return args.Get(0).(entities.Journey), args.Error(1)
}

func (m *JourneyMockUseCase) CheckJourney(ctx context.Context, userID int) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}
