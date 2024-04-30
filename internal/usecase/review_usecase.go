package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type QuestionUseCaseInterface interface {
	CreateReview(userID int, review entities.Review) error
	GetQuestions() ([]entities.QuestionResponse, error)
	CheckReview(userID int) (bool, error)
	SetStat(userID int) ([]entities.Statistic, error)
}

type QuestionUseCase struct {
	QuestionStorage storage.QuestionStorageInterface
}

func NewQuestionUseCase(storage storage.QuestionStorageInterface) QuestionUseCaseInterface {
	return &QuestionUseCase{
		QuestionStorage: storage,
	}
}

func (uc *QuestionUseCase) CreateReview(userID int, review entities.Review) error {
	return uc.QuestionStorage.AddReview(userID, review)
}

func (uc *QuestionUseCase) SetStat(userID int) ([]entities.Statistic, error) {
	AVGStat, err := uc.QuestionStorage.GetAvgStat()
	if err != nil {
		return []entities.Statistic{}, err
	}

	UserStat, err := uc.QuestionStorage.GetUserStat(userID)
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

func (uc *QuestionUseCase) GetQuestions() ([]entities.QuestionResponse, error) {
	return uc.QuestionStorage.GetQuestions()
}

func (uc *QuestionUseCase) CheckReview(userID int) (bool, error) {
	reviews, err := uc.QuestionStorage.GetReview(userID)
	if err != nil {
		return false, err
	}
	if len(reviews) == 0 {
		return false, nil
	}
	return true, nil
}
