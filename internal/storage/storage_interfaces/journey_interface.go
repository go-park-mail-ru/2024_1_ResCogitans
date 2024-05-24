package storage

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type JourneyStorageInterface interface {
	CreateJourney(ctx context.Context, journey entities.Journey) (entities.Journey, error)
	DeleteJourney(ctx context.Context, journeyID int) error
	GetJourneys(ctx context.Context, userID int) ([]entities.Journey, error)
	AddJourneySight(ctx context.Context, journeyID int, ids []int) error
	EditJourney(ctx context.Context, journeyID int, name, description string) error
	DeleteJourneySight(ctx context.Context, journeyID int, sight entities.JourneySight) error
	GetJourneySights(ctx context.Context, journeyID int) ([]*int, error)
	GetJourney(ctx context.Context, journeyID int) (entities.Journey, error)
}
