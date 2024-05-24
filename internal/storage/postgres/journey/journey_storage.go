package journey

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JourneyStorage struct {
	db *pgxpool.Pool
}

func NewJourneyStorage(db *pgxpool.Pool) *JourneyStorage {
	return &JourneyStorage{
		db: db,
	}
}

// CreateJourney создает новую поездку в базе данных.
func (js *JourneyStorage) CreateJourney(ctx context.Context, journey entities.Journey) (entities.Journey, error) {
	row := js.db.QueryRow(ctx,
		`INSERT INTO journey(name, user_id, description) 
			  VALUES ($1, $2, $3) RETURNING id, name, user_id, description;`, journey.Name, journey.UserID, journey.Description)
	err := row.Scan(&journey.ID, &journey.Name, &journey.UserID, &journey.Description)
	if err != nil {
		return entities.Journey{}, err
	}
	return journey, nil
}

func (js *JourneyStorage) DeleteJourney(ctx context.Context, journeyID int) error {
	_, err := js.db.Exec(ctx, `DELETE FROM journey_sight WHERE journey_id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	_, err = js.db.Exec(ctx, `DELETE FROM journey WHERE id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (js *JourneyStorage) GetJourneys(ctx context.Context, userID int) ([]entities.Journey, error) {
	var journey []*entities.Journey
	err := pgxscan.Select(ctx, js.db, &journey,
		`SELECT j.id, j.name, j.description, p.username FROM journey AS j 
    		   INNER JOIN profile_data AS p ON p.user_id = $1 
               WHERE j.user_id = $1;`, userID)
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
func (js *JourneyStorage) AddJourneySight(ctx context.Context, journeyID int, sightIDs []int) error {
	// Проверяем, существует ли journeyID в таблице journey
	var journeyExists bool
	err := js.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM journey WHERE id = $1)`, journeyID).Scan(&journeyExists)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	if !journeyExists {
		return fmt.Errorf("journey with id %d does not exist", journeyID)
	}

	// Добавляем достопримечательности в journey_sight
	for _, sightID := range sightIDs {
		_, err := js.db.Exec(ctx, `INSERT INTO journey_sight(journey_id, sight_id, priority) VALUES ($1, $2, $3)`, journeyID, sightID, 0)
		if err != nil {
			logger.Logger().Error(err.Error())
			return err
		}
	}

	return nil
}

func (js *JourneyStorage) EditJourney(ctx context.Context, journeyID int, name, description string) error {
	// Проверяем, существует ли journeyID в таблице journey
	var journeyExists bool
	err := js.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM journey WHERE id = $1)`, journeyID).Scan(&journeyExists)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	if !journeyExists {
		return fmt.Errorf("journey with id %d does not exist", journeyID)
	}

	// Обновляем имя и описание поездки
	_, err = js.db.Exec(ctx, `UPDATE journey SET name = $1, description = $2 WHERE id = $3`, name, description, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (js *JourneyStorage) DeleteJourneySight(ctx context.Context, journeyID int, sight entities.JourneySight) error {
	_, err := js.db.Exec(ctx, `DELETE FROM journey_sight WHERE journey_id = $1 AND sight_id = $2 `, journeyID, sight.SightID)
	return err
}

func (js *JourneyStorage) GetJourneySights(ctx context.Context, journeyID int) ([]*int, error) {
	var idList []*int
	err := pgxscan.Select(ctx, js.db, &idList, `SELECT js.sight_id FROM journey_sight AS js WHERE js.journey_id = $1`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}
	return idList, nil
}

func (js *JourneyStorage) GetJourney(ctx context.Context, journeyID int) (entities.Journey, error) {
	var journey []*entities.Journey

	err := pgxscan.Select(ctx, js.db, &journey,
		`SELECT j.id, j.name, j.description, p.username, p.user_id FROM journey AS j 
    			INNER JOIN profile_data AS p ON p.user_id = j.user_id 
                WHERE j.id = $1;`, journeyID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Journey{}, err
	}
	return *journey[0], nil
}
