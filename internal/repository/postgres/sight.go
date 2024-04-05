package repository

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/georgysavva/scany/v2/pgxscan"
)

// SightRepo struct
type SightRepo struct {
	db *pgxpool.Pool
}

// NewSightRepo creates sight repo
func NewSightRepo(db *pgxpool.Pool) *SightRepo {
	return &SightRepo{
		db: db,
	}
}

func (repo *SightRepo) GetSightsList() ([]entities.Sight, error) {
	var sights []*entities.Sight
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &sights, `SELECT sight.id, rating, name, description, city_id, country_id, image.path FROM sight INNER JOIN image ON sight.id = image.sight_id`)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	var sightList []entities.Sight
	for _, s := range sights {
		sightList = append(sightList, *s)
	}
	return sightList, nil
}

func (repo *SightRepo) GetSightByID(id int) (entities.Sight, error) {
	var sight []*entities.Sight
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &sight, `SELECT sight.id, rating, sight.name, description, city_id, country_id, image.path, city.city, country.country FROM sight INNER JOIN image ON sight.id = image.sight_id INNER JOIN city ON sight.city_id = city.id INNER JOIN country ON sight.country_id = country.id WHERE sight.id = $1`, id)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Sight{}, err
	}

	return *sight[0], nil
}

func (repo *SightRepo) GetCommentsBySightID(id int) ([]entities.Comment, error) {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &comments, `SELECT id, user_id, sight_id, rating, feedback FROM feedback WHERE sight_id = $1`, id)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	var commentsList []entities.Comment
	for _, s := range comments {
		commentsList = append(commentsList, *s)
	}

	return commentsList, nil
}

func (repo *SightRepo) CreateCommentBySightID(dataStr map[string]string, dataInt map[string]int) error {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &comments, `INSERT INTO feedback(user_id, sight_id, rating, feedback) VALUES($1, $2, $3, $4)`, dataInt["userID"], dataInt["sightID"], dataInt["rating"], dataStr["feedback"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *SightRepo) EditCommentByCommentID(dataStr map[string]string, dataInt map[string]int) error {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &comments, `UPDATE feedback SET rating = $1, feedback = $2 WHERE id = $3`, dataInt["rating"], dataStr["feedback"], dataInt["id"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *SightRepo) DeleteCommentByCommentID(dataInt map[string]int) error {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &comments, `DELETE FROM feedback WHERE id = $1`, dataInt["id"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}
