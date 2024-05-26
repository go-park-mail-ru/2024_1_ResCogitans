package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/comment"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
)

type SightUseCaseInterface interface {
	GetSightByID(ctx context.Context, sightID int) (entities.Sight, error)
	GetCommentsBySightID(ctx context.Context, commentID int) ([]entities.Comment, error)
	GetCommentsByUserID(ctx context.Context, userID int) ([]entities.Comment, error)
	GetSightsList(ctx context.Context) ([]entities.Sight, error)
	SearchSights(ctx context.Context, str string) (entities.Sights, error)
}

type SightUseCase struct {
	sightStorage   *sight.SightStorage
	commentStorage *comment.CommentStorage
}

func NewSightUseCase(sightStorage *sight.SightStorage, commentStorage *comment.CommentStorage) *SightUseCase {
	return &SightUseCase{
		sightStorage:   sightStorage,
		commentStorage: commentStorage,
	}
}

func (su *SightUseCase) GetSightByID(ctx context.Context, sightID int) (entities.Sight, error) {
	return su.sightStorage.GetSight(ctx, sightID)
}

func (su *SightUseCase) GetCommentsBySightID(ctx context.Context, commentID int) ([]entities.Comment, error) {
	return su.commentStorage.GetCommentsBySightID(ctx, commentID)
}

func (su *SightUseCase) GetCommentsByUserID(ctx context.Context, userID int) ([]entities.Comment, error) {
	return su.commentStorage.GetCommentsByUserID(ctx, userID)
}

func (su *SightUseCase) GetSightsList(ctx context.Context) ([]entities.Sight, error) {
	return su.sightStorage.GetSightsList(ctx)
}

func (su *SightUseCase) SearchSights(ctx context.Context, str string) (entities.Sights, error) {
	return su.sightStorage.SearchSights(ctx, str)
}
