package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/db"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DataBaseInitialization() (*pgxpool.Pool, *redis.Client, *redis.Client, error) {
	pdb, err := db.GetPostgres()
	if err != nil {
		return nil, nil, nil, err
	}

	rdb, err := db.GetRedis()
	if err != nil {
		pdb.Close()
		return nil, nil, nil, err
	}

	cdb, err := db.GetCSRFRedis()
	if err != nil {
		pdb.Close()
		_ = rdb.Close()
	}

	return pdb, rdb, cdb, nil
}
