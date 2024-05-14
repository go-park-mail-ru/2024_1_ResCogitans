package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
)

type CommentUseCaseInterface interface {
	CreateCommentBySightID(sightID int, comment entities.Comment) error
	EditCommentByCommentID(commentID int, comment entities.Comment) error
	DeleteCommentByCommentID(commentID int) error
	CheckCommentByUserID(userID int) (bool, error)
}

type CommentUseCase struct {
	SightStorage *sight.SightStorage
}

func NewCommentUseCase(storage *sight.SightStorage) *CommentUseCase {
	return &CommentUseCase{
		SightStorage: storage,
	}
}

func (cu *CommentUseCase) CreateCommentBySightID(sightID int, comment entities.Comment) error {
	return cu.SightStorage.CreateCommentBySightID(sightID, comment)
}

func (cu *CommentUseCase) EditCommentByCommentID(commentID int, comment entities.Comment) error {
	return cu.SightStorage.EditComment(commentID, comment)
}

func (cu *CommentUseCase) DeleteCommentByCommentID(commentID int) error {
	return cu.SightStorage.DeleteComment(commentID)
}

func (cu *CommentUseCase) CheckCommentByUserID(userID int) (bool, error) {
	comments, err := cu.SightStorage.GetCommentsByUserID(userID)
	if err != nil {
		return false, err
	}
	if len(comments) == 0 {
		return false, nil
	}
	return true, nil
}
