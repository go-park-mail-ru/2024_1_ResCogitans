package repository

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	su "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SightRepositoryI interface {
	NewSightRepo(db *pgxpool.Pool) *su.SightRepo
	GetSightsList() ([]entities.Sight, error)
	GetSightByID(id int) (entities.Sight, error)
	GetCommentsBySightID(id int) ([]entities.Comment, error)
}
