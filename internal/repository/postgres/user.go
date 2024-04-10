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

func (repo *UserRepo) GetUserProfile(dataStr map[string]int) (entities.UserProfile, error) {
	var user []entities.UserProfile
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &user, `SELECT user_id, username, bio, avatar FROM "profile" WHERE user_id = $1`, dataStr["userID"])

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
