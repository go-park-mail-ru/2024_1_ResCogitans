package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type CommentUseCaseInterface interface {
	CreateCommentBySightID(sightID int, comment entities.Comment) error
	EditCommentByCommentID(commentID int, comment entities.Comment) error
	DeleteCommentByCommentID(commentID int) error
}

type CommentUseCase struct {
	SightStorage storage.SightStorageInterface
}

func NewCommentUseCase(storage storage.SightStorageInterface) CommentUseCaseInterface {
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
