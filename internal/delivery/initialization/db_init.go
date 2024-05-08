package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/db"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DataBaseInitialization() (*pgxpool.Pool, *redis.Client, error) {
	pdb, err := db.GetPostgres()
	if err != nil {
		return nil, nil, err
	}

	rdb, err := db.GetRedis()
	if err != nil {
		pdb.Close()
		return nil, nil, err
	}

	return pdb, rdb, nil
}
