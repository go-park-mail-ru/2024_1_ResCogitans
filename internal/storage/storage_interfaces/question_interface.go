package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type QuestionInterface interface {
	AddReview(review entities.Review) error
	SetStat(userID int) ([]entities.Statistic, error)
}
