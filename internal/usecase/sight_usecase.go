package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type SightUseCaseInterface interface {
	GetSightByID(sightID int) (entities.Sight, error)
	GetCommentsBySightID(commentID int) ([]entities.Comment, error)
	GetCommentsByUserID(userID int) ([]entities.Comment, error)
	GetSightsList(ctx context.Context) ([]entities.Sight, error)
	SearchSights(searchParams map[string]string) (entities.Sights, error)
}

type SightUseCase struct {
	SightStorage storage.SightStorageInterface
}

func NewSightUseCase(storage storage.SightStorageInterface) SightUseCaseInterface {
	return &SightUseCase{
		SightStorage: storage,
	}
}

func (su *SightUseCase) GetSightByID(sightID int) (entities.Sight, error) {
	return su.SightStorage.GetSight(sightID)
}

func (su *SightUseCase) GetCommentsBySightID(commentID int) ([]entities.Comment, error) {
	return su.SightStorage.GetCommentsBySightID(commentID)
}

func (su *SightUseCase) GetCommentsByUserID(userID int) ([]entities.Comment, error) {
	return su.SightStorage.GetCommentsByUserID(userID)
}

func (su *SightUseCase) GetSightsList(ctx context.Context) ([]entities.Sight, error) {
	return su.SightStorage.GetSightsList(ctx)
}

func (su *SightUseCase) SearchSights(searchParams map[string]string) (entities.Sights, error) {
	return su.SightStorage.SearchSights(searchParams)
}
