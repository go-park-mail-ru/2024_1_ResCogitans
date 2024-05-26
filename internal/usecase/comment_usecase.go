package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/comment"
)

type CommentUseCaseInterface interface {
	CreateCommentBySightID(ctx context.Context, sightID int, comment entities.Comment) error
	EditCommentByCommentID(ctx context.Context, commentID int, comment entities.Comment) error
	DeleteCommentByCommentID(ctx context.Context, commentID int) error
	CheckCommentByUserID(ctx context.Context, userID int) (bool, error)
}

type CommentUseCase struct {
	storage *comment.CommentStorage
}

func NewCommentUseCase(storage *comment.CommentStorage) *CommentUseCase {
	return &CommentUseCase{
		storage: storage,
	}
}

func (cu *CommentUseCase) CreateCommentBySightID(ctx context.Context, sightID int, comment entities.Comment) error {
	return cu.storage.CreateCommentBySightID(ctx, sightID, comment)
}

func (cu *CommentUseCase) EditCommentByCommentID(ctx context.Context, commentID int, comment entities.Comment) error {
	return cu.storage.EditComment(ctx, commentID, comment)
}

func (cu *CommentUseCase) DeleteCommentByCommentID(ctx context.Context, commentID int) error {
	return cu.storage.DeleteComment(ctx, commentID)
}

func (cu *CommentUseCase) CheckCommentByUserID(ctx context.Context, userID int) (bool, error) {
	comments, err := cu.storage.GetCommentsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	if len(comments) == 0 {
		return false, nil
	}
	return true, nil
}
