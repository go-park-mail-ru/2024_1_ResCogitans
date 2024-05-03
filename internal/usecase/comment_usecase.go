package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type CommentUseCaseInterface interface {
	CreateCommentBySightID(sightID int, comment entities.CommentRequest) error
	EditCommentByCommentID(commentID int, comment entities.CommentRequest) error
	DeleteCommentByCommentID(commentID int) error
	CheckCommentByUserID(userID int) (bool, error)
}

type CommentUseCase struct {
	SightStorage storage.SightStorageInterface
}

func NewCommentUseCase(storage storage.SightStorageInterface) CommentUseCaseInterface {
	return &CommentUseCase{
		SightStorage: storage,
	}
}

func (cu *CommentUseCase) CreateCommentBySightID(sightID int, comment entities.CommentRequest) error {
	return cu.SightStorage.CreateCommentBySightID(sightID, comment)
}

func (cu *CommentUseCase) EditCommentByCommentID(commentID int, comment entities.CommentRequest) error {
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
