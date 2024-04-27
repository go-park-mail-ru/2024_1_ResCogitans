package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type QuestionUseCaseInterface interface {
	CreateReview(review entities.Review) error
}

type QuestionUseCase struct {
	QuestionStorage storage.QuestionInterface
}

func NewQuestionUseCase(storage storage.QuestionInterface) QuestionUseCaseInterface {
	return &QuestionUseCase{
		QuestionStorage: storage,
	}
}

func (uc *QuestionUseCase) CreateReview(review entities.Review) error {
	return uc.QuestionStorage.AddReview(review)
}
