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
