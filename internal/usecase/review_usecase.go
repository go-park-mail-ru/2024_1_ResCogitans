package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/question"
)

type QuestionUseCaseInterface interface {
	CreateReview(ctx context.Context, userID int, review entities.Review) error
	GetQuestions(ctx context.Context) ([]entities.QuestionResponse, error)
	CheckReview(ctx context.Context, userID int) (bool, error)
	SetStat(ctx context.Context, userID int) ([]entities.Statistic, error)
}

type QuestionUseCase struct {
	QuestionStorage *question.QuestionStorage
}

func NewQuestionUseCase(storage *question.QuestionStorage) *QuestionUseCase {
	return &QuestionUseCase{
		QuestionStorage: storage,
	}
}

func (uc *QuestionUseCase) CreateReview(ctx context.Context, userID int, review entities.Review) error {
	return uc.QuestionStorage.AddReview(ctx, userID, review)
}

func (uc *QuestionUseCase) SetStat(ctx context.Context, userID int) ([]entities.Statistic, error) {
	AVGStat, err := uc.QuestionStorage.GetAvgStat(ctx)
	if err != nil {
		return []entities.Statistic{}, err
	}

	UserStat, err := uc.QuestionStorage.GetUserStat(ctx, userID)
	if err != nil {
		return []entities.Statistic{}, err
	}

	if len(UserStat) == 0 {
		return AVGStat, nil
	}
	index := 0
	for i, question := range AVGStat {
		if question.ID == UserStat[index].ID {
			AVGStat[i].UserGrade = UserStat[index].UserGrade
			if index == len(UserStat)-1 {
				return AVGStat, nil
			}
			index++
		}
	}
	return AVGStat, nil
}

func (uc *QuestionUseCase) GetQuestions(ctx context.Context) ([]entities.QuestionResponse, error) {
	return uc.QuestionStorage.GetQuestions(ctx)
}

func (uc *QuestionUseCase) CheckReview(ctx context.Context, userID int) (bool, error) {
	reviews, err := uc.QuestionStorage.GetReview(ctx, userID)
	if err != nil {
		return false, err
	}
	if len(reviews) == 0 {
		return false, nil
	}
	return true, nil
}
