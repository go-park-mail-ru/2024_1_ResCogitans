package repository

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	errEditProfile = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed edit profile",
	}
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

	var userID int
	query := `INSERT INTO "user_data"(email, passwrd) VALUES ($1, $2) RETURNING id`
	err := repo.db.QueryRow(ctx, query, dataStr["email"], dataStr["passwrd"]).Scan(&userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.User{}, err
	}

	newUser := entities.User{
		ID:    userID,
		Email: dataStr["email"],
	}

	return newUser, nil
}

func (repo *UserRepo) AuthorizeUser(dataStr map[string]string) (entities.User, error) {
	var user []*entities.User
	ctx := context.Background()

	err := pgxscan.Select(ctx, repo.db, &user, `SELECT id, email, passwrd FROM user_data WHERE email = $1`, dataStr["email"])

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

	err := pgxscan.Select(ctx, repo.db, &user, `SELECT user_id, username, bio, avatar FROM profile_data WHERE user_id = $1`, dataInt["userID"])

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

	_, err := repo.db.Exec(ctx, `DELETE FROM profile_data WHERE user_id = $1`, dataInt["userID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}
	_, err = repo.db.Exec(ctx, `DELETE FROM user_data WHERE id = $1`, dataInt["userID"])
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (repo *UserRepo) EditUserProfile(dataInt map[string]int, dataStr map[string]string) (entities.UserProfile, error) {
	ctx := context.Background()

	var count int = 1

	query := "UPDATE profile_data SET "
	var queryParams []interface{}
	queryParams = append(queryParams, dataInt["userID"])
	var setClauses []string

	if dataStr["username"] != "" {
		count += 1
		str := "username = $" + strconv.Itoa(count)
		setClauses = append(setClauses, str)
		queryParams = append(queryParams, dataStr["username"])
	}
	if dataStr["bio"] != "" {
		count += 1
		str := "bio = $" + strconv.Itoa(count)
		setClauses = append(setClauses, str)
		queryParams = append(queryParams, dataStr["bio"])
	}
	if dataStr["avatar"] != "" {
		count += 1
		str := "avatar = $" + strconv.Itoa(count)
		setClauses = append(setClauses, str)
		queryParams = append(queryParams, dataStr["avatar"])
	}

	query += strings.Join(setClauses, ", ") + " WHERE user_id = $1"
	fmt.Println(query)

	_, err := repo.db.Exec(ctx, query, queryParams...)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.UserProfile{}, err
	}

	if count == 1 {
		logger.Logger().Error("No values to update")
		return entities.UserProfile{}, err
	}

	updatedProfile, err := repo.GetUserProfile(dataInt)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.UserProfile{}, errEditProfile
	}

	return updatedProfile, nil
}

func (repo *UserRepo) UpdateUserPassword(userID int, newPassword string) error {
	ctx := context.Background()
	logger := logger.Logger()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = repo.db.Exec(ctx, `UPDATE "user_data" SET passwrd = $1 WHERE id = $2`, string(hashedPassword), userID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
