package question

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionStorage struct {
	db *pgxpool.Pool
}

func NewQuestionStorage(db *pgxpool.Pool) *QuestionStorage {
	return &QuestionStorage{
		db: db,
	}
}

func (qs *QuestionStorage) AddReview(review entities.Review) error {
	currentTime := time.Now()
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	// Преобразуем текущее время в часовой пояс Москвы
	currentTimeInMoscow := currentTime.In(moscowLocation)
	timePlusThreeHours := currentTimeInMoscow.Add(3 * time.Hour)

	query := `
        INSERT INTO quiz (user_id, rating, question_id, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err = qs.db.Exec(context.Background(), query, review.UserID, review.Rating, review.QuestionID, timePlusThreeHours)
	return err
}

func (qs *QuestionStorage) GetQuestions() ([]entities.QuestionResponse, error) {
	var questions []*entities.QuestionResponse
	ctx := context.Background()
	err := pgxscan.Select(ctx, qs.db, &questions, `SELECT * FROM question`)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.QuestionResponse{}, err
	}

	var questionList []entities.QuestionResponse
	for _, q := range questions {
		questionList = append(questionList, *q)
	}

	return questionList, nil
}
