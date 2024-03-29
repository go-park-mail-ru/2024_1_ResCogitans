package db

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	_ "github.com/jackc/pgx/stdlib"
)

// GetPostgress gets connection to postgres database
func GetPostgres() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "mrdzhofik"
		password = "246858"
		dbname   = "jantugan"
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Logger().Error("Cannot open database")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.Logger().Error("Cannot connect to database")
		return nil, err
	}
	db.SetMaxOpenConns(10)

	return db, nil
}
