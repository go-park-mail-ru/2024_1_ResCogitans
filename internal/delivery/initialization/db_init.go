package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/db"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBs struct {
	PostgresDB *pgxpool.Pool
	CsrfDB     *redis.Client
}

// DataBaseInitialization Инициализация баз данных
func DataBaseInitialization(cfg *config.Config) (*DBs, error) {
	postgresDB, err := db.GetPostgres(cfg) // Инициализация основной базы данных
	if err != nil {
		return nil, err
	}

	csrfDB, err := db.GetCSRFRedis(cfg) // Инициализация базы данных для csrf
	if err != nil {
		postgresDB.Close()
		return nil, err
	}
	return &DBs{
		PostgresDB: postgresDB,
		CsrfDB:     csrfDB,
	}, nil
}
