package usecase

import (
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type CommentUseCaseInterface interface {
	CreateCommentBySightID(dataStr map[string]string, dataInt map[string]int) error
	EditCommentByCommentID(dataStr map[string]string, dataInt map[string]int) error
	DeleteCommentByCommentID(dataInt map[string]int) error
}

type CommentUseCase struct {
	SightStorage storage.SightStorageInterface
}

func NewCommentUseCase(storage storage.SightStorageInterface) CommentUseCaseInterface {
	return &CommentUseCase{
		SightStorage: storage,
	}
}

func (cu *CommentUseCase) CreateCommentBySightID(dataStr map[string]string, dataInt map[string]int) error {
	return cu.SightStorage.CreateCommentBySightID(dataStr, dataInt)
}

func (cu *CommentUseCase) EditCommentByCommentID(dataStr map[string]string, dataInt map[string]int) error {
	return cu.SightStorage.EditComment(dataStr, dataInt)
}

func (cu *CommentUseCase) DeleteCommentByCommentID(dataInt map[string]int) error {
	return cu.SightStorage.DeleteComment(dataInt)
}
