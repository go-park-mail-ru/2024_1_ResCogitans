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

	err := pgxscan.Select(ctx, ss.db, &sights, `SELECT sight.id, rating, name, description, city_id, country_id, image.path FROM sight INNER JOIN image ON sight.id = image.sight_id`)
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

func (ss *SightStorage) GetSight(sightID int) (entities.Sight, error) {
	var sight []*entities.Sight
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &sight, `SELECT sight.id, rating, sight.name, description, city_id, country_id, image.path, city.city, country.country FROM sight INNER JOIN image ON sight.id = image.sight_id INNER JOIN city ON sight.city_id = city.id INNER JOIN country ON sight.country_id = country.id WHERE sight.id = $1`, sightID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Sight{}, err
	}

	return *sight[0], nil
}

func (ss *SightStorage) GetCommentsBySightID(sightID int) ([]entities.Comment, error) {
	var comments []*entities.Comment
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &comments, `SELECT feedback.id, u.email, sight_id, rating, feedback FROM feedback INNER JOIN "user" AS u ON user_id = u.id WHERE sight_id = $1 `, sightID)
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

func (ss *SightStorage) CreateCommentBySightID(dataStr map[string]string, dataInt map[string]int) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `INSERT INTO feedback(user_id, sight_id, rating, feedback) VALUES($1, $2, $3, $4)`, dataInt["userID"], dataInt["sightID"], dataInt["rating"], dataStr["feedback"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) EditComment(dataStr map[string]string, dataInt map[string]int) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `UPDATE feedback SET rating = $1, feedback = $2 WHERE id = $3`, dataInt["rating"], dataStr["feedback"], dataInt["id"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) DeleteComment(dataInt map[string]int) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `DELETE FROM feedback WHERE id = $1`, dataInt["id"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) CreateJourney(dataInt map[string]int, dataStr map[string]string) (entities.Journey, error) {
	var journey entities.Journey
	ctx := context.Background()

	row := ss.db.QueryRow(ctx, `INSERT INTO journey(name, user_id, description) VALUES ($1, $2, $3) RETURNING id, name, user_id, description`, dataStr["name"], dataInt["userID"], dataStr["description"])
	err := row.Scan(&journey.ID, &journey.Name, &journey.UserID, &journey.Description)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Journey{}, err
	}

	return journey, nil
}

func (ss *SightStorage) DeleteJourney(dataInt map[string]int) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `DELETE FROM journey WHERE id = $1`, dataInt["journeyID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) GetJourneys(userID int) ([]entities.Journey, error) {
	var journey []*entities.Journey
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &journey, `SELECT j.id, j.name, j.description, u.email FROM journey AS j INNER JOIN "user" AS u ON u.id = $1`, userID)
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

func (ss *SightStorage) AddJourneySight(dataInt map[string]int) error {
	ctx := context.Background()
	var priority []*int
	var precedence int
	println(dataInt["sightID"])

	_ = pgxscan.Select(ctx, ss.db, &priority, `SELECT priority FROM journey_sight AS js WHERE js.journey_id = $1 ORDER BY priority DESC LIMIT 1;`, dataInt["journeyID"])
	if priority == nil {
		fmt.Println("Oops")
		precedence = 0
	} else {
		precedence = *priority[0]
	}

	_, err := ss.db.Exec(ctx, `INSERT INTO journey_sight(journey_id, sight_id, priority) VALUES($1, $2, $3) `, dataInt["journeyID"], dataInt["sightID"], precedence+1)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (ss *SightStorage) DeleteJourneySight(dataInt map[string]int) error {
	ctx := context.Background()

	_, err := ss.db.Exec(ctx, `DELETE FROM journey_sight WHERE journey_id = $1 AND sight_id = $2 `, dataInt["journeyID"], dataInt["sightID"])
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
			logger.Logger().Error(err.Error())
			continue
		}
		sights = append(sights, sight)
	}

	return sights, nil
}

func (ss *SightStorage) GetJourney(journeyID int) (entities.Journey, error) {
	var journey []*entities.Journey
	ctx := context.Background()

	err := pgxscan.Select(ctx, ss.db, &journey, `SELECT j.id, j.name, j.description, u.email FROM journey AS j INNER JOIN "user" AS u ON u.id = j.user_id WHERE j.id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Journey{}, err
	}

	return *journey[0], nil
}
