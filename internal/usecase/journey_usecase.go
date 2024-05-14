package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
)

type JourneyUseCaseInterface interface {
	CreateJourney(journey entities.Journey) (entities.Journey, error)
	DeleteJourneyByID(journeyID int) error
	GetJourneys(userID int) ([]entities.Journey, error)
	AddJourneySight(journeyID int, ids []int) error
	EditJourney(journeyID int, name, description string) error
	DeleteJourneySight(journeyID int, sight entities.JourneySight) error
	GetJourneySights(journeyID int) ([]entities.Sight, error)
	GetJourney(journeyID int) (entities.Journey, error)
	CheckJourney(userID int) (bool, error)
}

type JourneyUseCase struct {
	SightStorage *sight.SightStorage
}

func NewJourneyUseCase(storage *sight.SightStorage) *JourneyUseCase {
	return &JourneyUseCase{
		SightStorage: storage,
	}
}

func (ju *JourneyUseCase) CreateJourney(journey entities.Journey) (entities.Journey, error) {
	return ju.SightStorage.CreateJourney(journey)
}

func (ju *JourneyUseCase) DeleteJourneyByID(journeyID int) error {
	return ju.SightStorage.DeleteJourney(journeyID)
}

func (ju *JourneyUseCase) GetJourneys(userID int) ([]entities.Journey, error) {
	return ju.SightStorage.GetJourneys(userID)
}

func (ju *JourneyUseCase) AddJourneySight(journeyID int, ids []int) error {
	return ju.SightStorage.AddJourneySight(journeyID, ids)
}

func (ju *JourneyUseCase) EditJourney(journeyID int, name, description string) error {
	return ju.SightStorage.EditJourney(journeyID, name, description)
}

func (ju *JourneyUseCase) DeleteJourneySight(journeyID int, sight entities.JourneySight) error {
	return ju.SightStorage.DeleteJourneySight(journeyID, sight)
}

func (ju *JourneyUseCase) GetJourneySights(journeyID int) ([]entities.Sight, error) {
	return ju.SightStorage.GetJourneySights(journeyID)
}

func (ju *JourneyUseCase) GetJourney(journeyID int) (entities.Journey, error) {
	return ju.SightStorage.GetJourney(journeyID)
}

func (ju *JourneyUseCase) CheckJourney(userID int) (bool, error) {
	journeys, err := ju.SightStorage.GetJourneys(userID)
	if err != nil {
		return false, err
	}
	if len(journeys) == 0 {
		return false, nil
	}
	return true, nil
}
