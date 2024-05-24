package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/journey"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
)

type JourneyUseCaseInterface interface {
	CreateJourney(ctx context.Context, journey entities.Journey) (entities.Journey, error)
	DeleteJourneyByID(ctx context.Context, journeyID int) error
	GetJourneys(ctx context.Context, userID int) ([]entities.Journey, error)
	AddJourneySight(ctx context.Context, journeyID int, ids []int) error
	EditJourney(ctx context.Context, journeyID int, name, description string) error
	DeleteJourneySight(ctx context.Context, journeyID int, sight entities.JourneySight) error
	GetJourneySights(ctx context.Context, journeyID int) ([]entities.Sight, error)
	GetJourney(ctx context.Context, journeyID int) (entities.Journey, error)
	CheckJourney(ctx context.Context, userID int) (bool, error)
}

type JourneyUseCase struct {
	journeyStorage *journey.JourneyStorage
	sightStorage   *sight.SightStorage
}

func NewJourneyUseCase(storage *journey.JourneyStorage) *JourneyUseCase {
	return &JourneyUseCase{
		journeyStorage: storage,
	}
}

func (ju *JourneyUseCase) CreateJourney(ctx context.Context, journey entities.Journey) (entities.Journey, error) {
	return ju.journeyStorage.CreateJourney(ctx, journey)
}

func (ju *JourneyUseCase) DeleteJourneyByID(ctx context.Context, journeyID int) error {
	return ju.journeyStorage.DeleteJourney(ctx, journeyID)
}

func (ju *JourneyUseCase) GetJourneys(ctx context.Context, userID int) ([]entities.Journey, error) {
	return ju.journeyStorage.GetJourneys(ctx, userID)
}

func (ju *JourneyUseCase) AddJourneySight(ctx context.Context, journeyID int, ids []int) error {
	return ju.journeyStorage.AddJourneySight(ctx, journeyID, ids)
}

func (ju *JourneyUseCase) EditJourney(ctx context.Context, journeyID int, name, description string) error {
	return ju.journeyStorage.EditJourney(ctx, journeyID, name, description)
}

func (ju *JourneyUseCase) DeleteJourneySight(ctx context.Context, journeyID int, sight entities.JourneySight) error {
	return ju.journeyStorage.DeleteJourneySight(ctx, journeyID, sight)
}

func (ju *JourneyUseCase) GetJourneySights(ctx context.Context, journeyID int) ([]entities.Sight, error) {
	idList, err := ju.journeyStorage.GetJourneySights(ctx, journeyID)
	if err != nil {
		return []entities.Sight{}, err
	}
	var sights []entities.Sight
	for _, id := range idList {
		getSight, err := ju.sightStorage.GetSight(ctx, *id)
		if err != nil {
			return nil, err
		}
		sights = append(sights, getSight)
	}

	return sights, nil
}

func (ju *JourneyUseCase) GetJourney(ctx context.Context, journeyID int) (entities.Journey, error) {
	return ju.journeyStorage.GetJourney(ctx, journeyID)
}

func (ju *JourneyUseCase) CheckJourney(ctx context.Context, userID int) (bool, error) {
	journeys, err := ju.journeyStorage.GetJourneys(ctx, userID)
	if err != nil {
		return false, err
	}
	if len(journeys) == 0 {
		return false, nil
	}
	return true, nil
}
