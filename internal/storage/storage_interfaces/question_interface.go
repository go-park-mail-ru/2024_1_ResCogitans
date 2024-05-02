package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type QuestionStorageInterface interface {
	AddReview(userID int, review entities.ReviewRequest) error
	SetStat(userID int) ([]entities.Statistic, error)
	GetQuestions() ([]entities.QuestionResponse, error)
	GetReview(userID int) ([]entities.Review, error)
	GetAvgStat() ([]entities.Statistic, error)
	GetUserStat(userID int) ([]entities.Statistic, error)
}
