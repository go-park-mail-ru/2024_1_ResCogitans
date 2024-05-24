package sight

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/georgysavva/scany/v2/pgxscan"
)

// SightStorage struct
type SightStorage struct {
	db *pgxpool.Pool
}

// NewSightStorage creates sight repo
func NewSightStorage(db *pgxpool.Pool) *SightStorage {
	return &SightStorage{
		db: db,
	}
}

func (ss *SightStorage) GetSightsList(ctx context.Context) ([]entities.Sight, error) {
	var sights []*entities.Sight

	err := pgxscan.Select(ctx, ss.db, &sights, `SELECT sight.id, COALESCE(rating, 0) AS rating, name, description, city_id, country_id, im.path FROM sight 
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

func (ss *SightStorage) GetSight(ctx context.Context, id int) (entities.Sight, error) {
	var sight []*entities.Sight

	query := `SELECT sight.id, COALESCE(rating, 0) AS rating, sight.name, description, city_id, sight.country_id, im.path, city.city, country.country 
	FROM sight 
	INNER JOIN image_data AS im 
	ON sight.id = im.sight_id 
	INNER JOIN city 
	ON sight.city_id = city.id 
	INNER JOIN country 
	ON sight.country_id = country.id 
	WHERE sight.id = $1`

	err := pgxscan.Select(ctx, ss.db, &sight, query, id)
	if err != nil {
		return entities.Sight{}, err
	}

	return *sight[0], nil
}

func (ss *SightStorage) SearchSights(ctx context.Context, str string) (entities.Sights, error) {
	query := `SELECT id, rating, name, description, city_id, country_id
    FROM sight
    WHERE LOWER(name) LIKE LOWER($1)`

	rows, err := ss.db.Query(ctx, query, "%"+str+"%")
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
