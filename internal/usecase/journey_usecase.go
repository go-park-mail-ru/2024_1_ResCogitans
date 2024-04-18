package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type JourneyUseCaseInterface interface {
	CreateJourney(dataInt map[string]int, dataStr map[string]string) (entities.Journey, error)
	DeleteJourneyByID(dataInt map[string]int) error
	GetJourneys(userID int) ([]entities.Journey, error)
	AddJourneySight(dataInt map[string]int) error
	DeleteJourneySight(dataInt map[string]int) error
	GetJourneySights(journeyID int) ([]entities.Sight, error)
	GetJourney(journeyID int) (entities.Journey, error)
}

type JourneyUseCase struct {
	SightStorage storage.SightStorageInterface
}

func NewJourneyUseCase(storage storage.SightStorageInterface) JourneyUseCaseInterface {
	return &JourneyUseCase{
		SightStorage: storage,
	}
}

func (ju *JourneyUseCase) CreateJourney(dataInt map[string]int, dataStr map[string]string) (entities.Journey, error) {
	return ju.SightStorage.CreateJourney(dataInt, dataStr)
}

func (ju *JourneyUseCase) DeleteJourneyByID(dataInt map[string]int) error {
	return ju.SightStorage.DeleteJourney(dataInt)
}

func (ju *JourneyUseCase) GetJourneys(userID int) ([]entities.Journey, error) {
	return ju.SightStorage.GetJourneys(userID)
}

func (ju *JourneyUseCase) AddJourneySight(dataInt map[string]int) error {
	return ju.SightStorage.AddJourneySight(dataInt)
}

func (ju *JourneyUseCase) DeleteJourneySight(dataInt map[string]int) error {
	return ju.SightStorage.DeleteJourneySight(dataInt)
}

func (ju *JourneyUseCase) GetJourneySights(journeyID int) ([]entities.Sight, error) {
	return ju.SightStorage.GetJourneySights(journeyID)
}

func (ju *JourneyUseCase) GetJourney(journeyID int) (entities.Journey, error) {
	return ju.SightStorage.GetJourney(journeyID)
}
