package user

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// UserStorage struct
type UserStorage struct {
	db *pgxpool.Pool
}

// NewUserStorage creates postgres storage
func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (us *UserStorage) SaveUser(ctx context.Context, email, password, salt string) error {
	_, err := us.db.Exec(ctx, `INSERT INTO user_data (email, passwrd, salt) VALUES ($1, $2, $3)`, email, password, salt)
	return err
}

func (us *UserStorage) ChangeEmail(ctx context.Context, userID int, email string) error {
	_, err := us.db.Exec(ctx, `UPDATE user_data SET email = $1 WHERE id = $2`, email, userID)
	return err
}

func (us *UserStorage) ChangePassword(ctx context.Context, userID int, password, salt string) error {
	_, err := us.db.Exec(ctx, `UPDATE user_data SET passwrd = $1, salt = $2 WHERE id = $3`, password, salt, userID)
	return err
}

func (us *UserStorage) GetUserByID(ctx context.Context, id int) (entities.User, error) {
	var user entities.User
	err := us.db.QueryRow(ctx, `SELECT id, email, passwrd, salt FROM user_data WHERE id = $1`, id).Scan(&user.ID, &user.Username, &user.Password, &user.Salt)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (us *UserStorage) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User
	err := us.db.QueryRow(ctx, `SELECT id, email, passwrd, salt FROM user_data WHERE email = $1`, email).Scan(&user.ID, &user.Username, &user.Password, &user.Salt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, httperrors.NewHttpError(http.StatusBadRequest, "user not found")
		}
		return entities.User{}, err
	}
	return user, nil
}

func (us *UserStorage) DeleteUser(ctx context.Context, userID int) error {
	_, err := us.db.Exec(ctx, `DELETE FROM user_data WHERE id = $1`, userID)
	return err
}

func (us *UserStorage) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var count int
	err := us.db.QueryRow(ctx, `SELECT COUNT(*) FROM user_data WHERE email = $1`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
