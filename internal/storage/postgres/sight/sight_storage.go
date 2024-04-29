package sight

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/georgysavva/scany/v2/pgxscan"
)

// SightStorage struct
type SightStorage struct {
	db *pgxpool.Pool
}

// NewSightRepo creates sight repo
func NewSightStorage(db *pgxpool.Pool) storage.SightStorageInterface {
	return &SightStorage{
		db: db,
	}
}

func (ss *SightStorage) GetSightsList() ([]entities.Sight, error) {
	var sights []*entities.Sight
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &sights, `SELECT sight.id, COALESCE(rating, 0) AS rating, name, description, city_id, country_id, im.path 
	FROM sight 
	INNER JOIN image_data AS im 
		ON sight.id = im.sight_id `)
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

func (ss *SightStorage) GetSight(id int) (entities.Sight, error) {
	var sight []*entities.Sight
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &sight, `SELECT sight.id, COALESCE(rating, 0) AS rating, sight.name, description, city_id, sight.country_id, im.path, city.city, country.country, longitude, latitude 
	FROM sight 
	INNER JOIN image_data AS im 
		ON sight.id = im.sight_id 
	INNER JOIN city 
		ON sight.city_id = city.id 
	INNER JOIN country 
		ON sight.country_id = country.id 
	WHERE sight.id = $1`, id)
	if err != nil {
		return entities.Sight{}, err
	}

	return *sight[0], nil
}

func (ss *SightStorage) SearchSights(str string) (entities.Sights, error) {
	query := `
        SELECT id, rating, name, description, city_id, country_id
        FROM sight
        WHERE LOWER(name) LIKE LOWER($1)
    `
	rows, err := ss.db.Query(context.Background(), query, "%"+str+"%")
	if err != nil {
		return entities.Sights{}, err
	}
	defer rows.Close()

	var sights entities.Sights
	for rows.Next() {
		var sight entities.Sight
		err := rows.Scan(&sight.ID, &sight.Rating, &sight.Name, &sight.Description, &sight.CityID, &sight.CountryID)
		if err != nil {
			return entities.Sights{}, err
		}
		sights.Sight = append(sights.Sight, sight)
	}

	if err := rows.Err(); err != nil {
		return entities.Sights{}, err
	}

	return sights, nil
}

func (ss *SightStorage) GetCommentsBySightID(id int) ([]entities.Comment, error) {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &comments, `SELECT f.id, f.user_id, p.username, p.avatar, f.sight_id, f.rating, f.feedback FROM feedback AS f INNER JOIN profile_data AS p ON f.user_id = p.user_id WHERE sight_id =  $1 `, id)
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

func (ss *SightStorage) GetCommentsByUserID(userID int) ([]entities.Comment, error) {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &comments, `SELECT f.id, f.user_id, f.sight_id, f.rating, f.feedback FROM feedback AS f WHERE user_id =  $1 `, userID)
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

func (ss *SightStorage) CreateCommentBySightID(sightID int, comment entities.Comment) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `INSERT INTO feedback(user_id, sight_id, rating, feedback) VALUES($1, $2, $3, $4)`, comment.UserID, sightID, comment.Rating, comment.Feedback)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) EditComment(commentID int, comment entities.Comment) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `UPDATE feedback SET rating = $1, feedback = $2 WHERE id = $3`, comment.Rating, comment.Feedback, commentID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) DeleteComment(commentID int) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `DELETE FROM feedback WHERE id = $1`, commentID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

// CreateJourney создает новую поездку в базе данных.
func (ss *SightStorage) CreateJourney(journey entities.Journey) (entities.Journey, error) {
	ctx := context.Background()

	row := ss.db.QueryRow(ctx, `INSERT INTO journey(name, user_id, description) VALUES ($1, $2, $3) RETURNING id, name, user_id, description;`, journey.Name, journey.UserID, journey.Description)
	err := row.Scan(&journey.ID, &journey.Name, &journey.UserID, &journey.Description)
	if err != nil {
		return entities.Journey{}, err
	}
	return journey, nil
}

func (ss *SightStorage) DeleteJourney(journeyID int) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `DELETE FROM journey_sight WHERE journey_id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	_, err = ss.db.Exec(ctx, `DELETE FROM journey WHERE id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) GetJourneys(userID int) ([]entities.Journey, error) {
	var journey []*entities.Journey
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &journey, `SELECT j.id, j.name, j.description, p.username FROM journey AS j INNER JOIN profile_data AS p ON p.user_id = $1 WHERE j.user_id = $1;`, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	var journeyList []entities.Journey
	for _, j := range journey {
		journeyList = append(journeyList, *j)
	}
	return journeyList, nil
}

// AddJourneySight добавляет достопримечательности в существующую поездку.
func (ss *SightStorage) AddJourneySight(journeyID int, sightIDs []int) error {
	ctx := context.Background()

	// Проверяем, существует ли journeyID в таблице journey
	var journeyExists bool
	err := ss.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM journey WHERE id = $1)`, journeyID).Scan(&journeyExists)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	if !journeyExists {
		return fmt.Errorf("journey with id %d does not exist", journeyID)
	}

	// Добавляем достопримечательности в journey_sight
	for _, sightID := range sightIDs {
		_, err := ss.db.Exec(ctx, `INSERT INTO journey_sight(journey_id, sight_id, priority) VALUES ($1, $2, $3)`, journeyID, sightID, 0)
		if err != nil {
			logger.Logger().Error(err.Error())
			return err
		}
	}

	return nil
}

func (ss *SightStorage) EditJourney(journeyID int, name, description string) error {
	ctx := context.Background()

	// Проверяем, существует ли journeyID в таблице journey
	var journeyExists bool
	err := ss.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM journey WHERE id = $1)`, journeyID).Scan(&journeyExists)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	if !journeyExists {
		return fmt.Errorf("journey with id %d does not exist", journeyID)
	}

	// Обновляем имя и описание поездки
	_, err = ss.db.Exec(ctx, `UPDATE journey SET name = $1, description = $2 WHERE id = $3`, name, description, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) GetJourneySights(journeyID int) ([]entities.Sight, error) {
	var sights []entities.Sight
	var idList []*int
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &idList, `SELECT js.sight_id FROM journey_sight AS js WHERE js.journey_id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	for _, id := range idList {
		sight, err := ss.GetSight(*id)
		if err != nil {
			return nil, err
		}
		sights = append(sights, sight)
	}

	return sights, nil
}

func (ss *SightStorage) GetJourney(journeyID int) (entities.Journey, error) {
	var journey []*entities.Journey
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &journey, `SELECT j.id, j.name, j.description, p.username, p.user_id FROM journey AS j INNER JOIN profile_data AS p ON p.user_id = j.user_id WHERE j.id = $1;`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Journey{}, err
	}
	return *journey[0], nil
}
