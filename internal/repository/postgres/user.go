package repository

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SightRepo struct
type UserRepo struct {
	db *pgxpool.Pool
}

// NewSightRepo creates sight repo
func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (repo *UserRepo) CreateUser(dataStr map[string]string) (entities.User, error) {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx, `INSERT INTO "user"(email, passwrd) VALUES ($1, $2)`, dataStr["email"], dataStr["passwrd"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.User{}, err
	}

	return entities.User{}, nil
}
