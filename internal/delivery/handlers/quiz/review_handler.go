package quiz

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
)

type QuizHandler struct {
	questionUseCase usecase.QuestionUseCaseInterface
}

func NewQuizHandler(questionUseCase usecase.QuestionUseCaseInterface) *QuizHandler {
	return &QuizHandler{
		questionUseCase: questionUseCase,
	}
}

func (h *QuizHandler) CreateReview(_ context.Context, requestData entities.Review) (bool, error) {
	err := h.questionUseCase.CreateReview(requestData)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *QuizHandler) CheckData(ctx context.Context, _ entities.Review) (bool, error) {}
