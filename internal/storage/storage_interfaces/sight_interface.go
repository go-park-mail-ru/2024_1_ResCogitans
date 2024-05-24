package storage

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type SightStorageInterface interface {
	GetSightsList(ctx context.Context) ([]entities.Sight, error)
	GetSight(ctx context.Context, sightID int) (entities.Sight, error)
	SearchSights(ctx context.Context, str string) (entities.Sights, error)
}
