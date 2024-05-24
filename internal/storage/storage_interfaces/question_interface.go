package storage

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type QuestionInterface interface {
	AddReview(ctx context.Context, userID int, review entities.Review) error
	SetStat(ctx context.Context, userID int) ([]entities.Statistic, error)
	GetQuestions(ctx context.Context) ([]entities.QuestionResponse, error)
	GetReview(ctx context.Context, userID int) ([]entities.Review, error)
	GetAvgStat(ctx context.Context) ([]entities.Statistic, error)
	GetUserStat(ctx context.Context, userID int) ([]entities.Statistic, error)
}
