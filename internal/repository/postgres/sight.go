package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

// SightRepo struct
type SightRepo struct {
	DB *sql.DB
}

// NewSighttRepo creates product repo
func NewSightRepo(db *sql.DB, err error) *SightRepo {
	return &SightRepo{
		DB: db,
	}
}

func (repo *SightRepo) GetSightsList() ([]*entities.Sight, error) {
	rows, err := repo.DB.Query(`SELECT sight.id, rating, name, description, city_id, country_id, image.path FROM sight INNER JOIN image ON sight.id = image.sight_id`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Sights = []*entities.Sight{}
	for rows.Next() {
		sight := &entities.Sight{}
		err = rows.Scan(
			&sight.ID,
			&sight.Rating,
			&sight.Name,
			&sight.Description,
			&sight.City_id,
			&sight.Country_id,
			&sight.Url,
		)
		if err != nil {
			return nil, err
		}
		Sights = append(Sights, sight)
	}
	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return Sights, nil
}
