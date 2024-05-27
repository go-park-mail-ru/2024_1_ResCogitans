package storage

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type CommentStorageInterface interface {
	GetCommentsBySightID(ctx context.Context, commentID int) ([]entities.Comment, error)
	GetCommentsByUserID(ctx context.Context, userID int) ([]entities.Comment, error)
	CreateCommentBySightID(ctx context.Context, sightID int, comment entities.Comment) error
	EditComment(ctx context.Context, commentID int, comment entities.Comment) error
	DeleteComment(ctx context.Context, commentID int) error
}
