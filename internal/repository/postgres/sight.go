package repository

import (
	"context"
	"fmt"

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

	err := pgxscan.Select(ctx, repo.db, &sights, `SELECT sight.id, COALESCE(rating, 0) AS rating, name, description, city_id, country_id, im.path 
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

func (repo *SightRepo) GetSightByID(id int) (entities.Sight, error) {
	var sight []*entities.Sight
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &sight, `SELECT sight.id, COALESCE(rating, 0) AS rating, sight.name, description, city_id, sight.country_id, im.path, city.city, country.country 
	FROM sight 
	INNER JOIN image_data AS im 
		ON sight.id = im.sight_id 
	INNER JOIN city 
		ON sight.city_id = city.id 
	INNER JOIN country 
		ON sight.country_id = country.id 
	WHERE sight.id = $1`, id)
	fmt.Println(id)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Sight{}, err
	}

	return *sight[0], nil
}

func (repo *SightRepo) GetCommentsBySightID(id int) ([]entities.Comment, error) {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &comments, `SELECT f.id, f.user_id, p.username, p.avatar, f.sight_id, f.rating, f.feedback FROM feedback AS f INNER JOIN profile_data AS p ON f.user_id = p.user_id WHERE sight_id =  $1 `, id)
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
	ctx := context.Background()

	_, err := repo.db.Exec(ctx, `INSERT INTO feedback(user_id, sight_id, rating, feedback) VALUES($1, $2, $3, $4)`, dataInt["userID"], dataInt["sightID"], dataInt["rating"], dataStr["feedback"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *SightRepo) EditCommentByCommentID(dataStr map[string]string, dataInt map[string]int) error {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx, `UPDATE feedback SET rating = $1, feedback = $2 WHERE id = $3`, dataInt["rating"], dataStr["feedback"], dataInt["id"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *SightRepo) DeleteCommentByCommentID(dataInt map[string]int) error {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx, `DELETE FROM feedback WHERE id = $1`, dataInt["id"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *SightRepo) CreateJourney(dataInt map[string]int, dataStr map[string]string) (entities.Journey, error) {
	var journey entities.Journey
	ctx := context.Background()

	row := repo.db.QueryRow(ctx, `INSERT INTO journey(name, user_id, description) VALUES ($1, $2, $3) RETURNING id, name, user_id, description;`, dataStr["name"], dataInt["userID"], dataStr["description"])
	err := row.Scan(&journey.ID, &journey.Name, &journey.UserID, &journey.Description)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Journey{}, err
	}

	return journey, nil
}

func (repo *SightRepo) DeleteJourneyByID(dataInt map[string]int) error {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx, `DELETE FROM journey_sight WHERE journey_id = $1`, dataInt["journeyID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	_, err = repo.db.Exec(ctx, `DELETE FROM journey WHERE id = $1`, dataInt["journeyID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *SightRepo) GetJourneys(userID int) ([]entities.Journey, error) {
	var journey []*entities.Journey
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &journey, `SELECT j.id, j.name, j.description, p.username FROM journey AS j INNER JOIN profile_data AS p ON p.user_id = $1 WHERE j.user_id = $1`, userID)
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

func (repo *SightRepo) AddJourneySight(dataInt map[string]int, ids []int, dataStr map[string]string) error {
	ctx := context.Background()
	precedence := 0

	_, err := repo.db.Exec(ctx, `UPDATE journey SET name = $1, description = $2 WHERE id = $3;`, dataStr["name"], dataStr["description"], dataInt["journeyID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}
	fmt.Println("Succes update")

	_, err = repo.db.Exec(ctx, `DELETE FROM journey_sight WHERE journey_sight.journey_id = $1;`, dataInt["journeyID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	for _, id := range ids {
		precedence += 1
		_, err := repo.db.Exec(ctx, `INSERT INTO journey_sight(journey_id, sight_id, priority) VALUES($1, $2, $3) `, dataInt["journeyID"], id, precedence)
		if err != nil {
			logger.Logger().Error(err.Error())
			return err
		}
	}

	return nil
}

func (repo *SightRepo) DeleteJourneySight(dataInt map[string]int) error {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx, `DELETE FROM journey_sight WHERE journey_id = $1 AND sight_id = $2 `, dataInt["journeyID"], dataInt["sightID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *SightRepo) GetJourneySights(journeyID int) ([]entities.Sight, error) {
	var sights []entities.Sight
	var idList []*int
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &idList, `SELECT js.sight_id FROM journey_sight AS js WHERE js.journey_id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	for _, id := range idList {
		sight, err := repo.GetSightByID(*id)
		if err != nil {
			logger.Logger().Error(err.Error())
			continue
		}
		sights = append(sights, sight)
	}

	return sights, nil
}

func (repo *SightRepo) GetJourney(journeyID int) (entities.Journey, error) {
	var journey []*entities.Journey
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &journey, `SELECT j.id, j.name, j.description, p.username, p.user_id FROM journey AS j INNER JOIN profile_data AS p ON p.user_id = j.user_id WHERE j.id = $1;`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Journey{}, err
	}

	fmt.Println(*journey[0])

	return *journey[0], nil
}
