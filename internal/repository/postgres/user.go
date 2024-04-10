package repository

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo struct
type UserRepo struct {
	db *pgxpool.Pool
}

// NewUserRepo creates sight repo
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

func (repo *UserRepo) AuthorizeUser(dataStr map[string]string) (entities.User, error) {
	var user []*entities.User
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &user, `SELECT id, email, passwrd FROM "user" WHERE email = $1`, dataStr["email"])

	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.User{}, err
	}

	if user == nil {
		fmt.Println("Empty user")
		return entities.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user[0].Passwrd), []byte(dataStr["passwrd"]))
	if user == nil {
		fmt.Println("Passwords not match!")
		return entities.User{}, err
	}

	return *user[0], nil
}

func (repo *UserRepo) GetUserProfile(dataInt map[string]int) (entities.UserProfile, error) {
	var user []entities.UserProfile
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &user, `SELECT user_id, username, bio, avatar FROM "profile" WHERE user_id = $1`, dataInt["userID"])

	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.UserProfile{}, err
	}

	if user == nil {
		fmt.Println("Empty user")
		return entities.UserProfile{}, err
	}

	return user[0], nil
}

func (repo *UserRepo) DeleteUserProfile(dataInt map[string]int) error {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx, `DELETE FROM "profile" WHERE user_id = $1`, dataInt["userID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}
	_, err = repo.db.Exec(ctx, `DELETE FROM "user" WHERE id = $1`, dataInt["userID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *UserRepo) EditUserProfile(dataInt map[string]int, dataStr map[string]string) (entities.UserProfile, error) {
	var profile entities.UserProfile
	ctx := context.Background()

	query := "UPDATE journey SET "
	var queryParams []interface{}

	// Добавляем поля для обновления в запрос
	if dataStr["name"] != "" {
		query += "name = $1, "
		queryParams = append(queryParams, dataStr["name"])
	}
	if dataStr["bio"] != "" {
		query += "bio = $2, "
		queryParams = append(queryParams, dataStr["bio"])
	}
	if dataStr["avatar"] != "" {
		query += "avatar = $3"
		queryParams = append(queryParams, dataStr["avatar"])
	}

	// Убираем последнюю запятую и добавляем условие WHERE
	query = query + " WHERE id = $4 RETURNING id, name, user_id, description"
	queryParams = append(queryParams, dataInt["userID"])

	row := repo.db.QueryRow(ctx, query, queryParams...)
	err := row.Scan(&profile.UserID, &profile.Username, &profile.Bio, &profile.Avatar)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.UserProfile{}, err
	}

	return profile, nil
}
