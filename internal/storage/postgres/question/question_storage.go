package question

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionStorage struct {
	db *pgxpool.Pool
}

func NewQuestionStorage(db *pgxpool.Pool) storage.QuestionStorageInterface {
	return &QuestionStorage{
		db: db,
	}
}

func (qs *QuestionStorage) AddReview(userID int, review entities.ReviewRequest) error {
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
	_, err = qs.db.Exec(context.Background(), query, userID, review.Rating, review.QuestionID, timePlusThreeHours)
	return err
}

func (qs *QuestionStorage) SetStat(userID int) ([]entities.Statistic, error) {
	var statistic []*entities.Statistic
	ctx := context.Background()
	err := pgxscan.Select(ctx, qs.db, &statistic, `SELECT q.text, r.rating AS user_grade, AVG(r.rating) AS average_grade
	FROM quiz r  
	INNER JOIN question q ON r.question_id = q.id 
-- 	WHERE user_id = $1 
	GROUP BY r.question_id, q.text, r.rating`, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.Statistic{}, err
	}

	var statisticList []entities.Statistic
	for _, s := range statistic {
		statisticList = append(statisticList, *s)
	}

	return statisticList, nil
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

func (qs *QuestionStorage) GetReview(userID int) ([]entities.Review, error) {
	var review []*entities.Review
	ctx := context.Background()
	err := pgxscan.Select(ctx, qs.db, &review, `SELECT * FROM quiz WHERE user_id = $1`, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.Review{}, err
	}

	var reviewList []entities.Review
	for _, r := range review {
		reviewList = append(reviewList, *r)
	}

	return reviewList, nil
}

func (qs *QuestionStorage) GetAvgStat() ([]entities.Statistic, error) {
	var statistic []*entities.Statistic
	ctx := context.Background()
	err := pgxscan.Select(ctx, qs.db, &statistic, `SELECT q.id, q.text, AVG(r.rating),
	FROM quiz r 
	INNER JOIN question q ON r.question_id = q.id 
	GROUP BY r.question_id
	ORDER BY q.id`)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.Statistic{}, err
	}

	var statisticList []entities.Statistic
	for _, s := range statistic {
		statisticList = append(statisticList, *s)
	}

	return statisticList, nil
}

func (qs *QuestionStorage) GetUserStat(userID int) ([]entities.Statistic, error) {
	var statistic []*entities.Statistic
	ctx := context.Background()

	err := pgxscan.Select(ctx, qs.db, &statistic, `SELECT q.id, q.text, r.rating 
	FROM question q 
	INNER JOIN quiz r ON q.id = r.question_id 
	WHERE user_id = $1
	ORDER BY q.id `, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.Statistic{}, err
	}

	var statisticList []entities.Statistic
	for _, s := range statistic {
		statisticList = append(statisticList, *s)
	}

	return statisticList, nil
}
