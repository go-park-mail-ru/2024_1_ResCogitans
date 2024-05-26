package user

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// UserProfileStorage struct
type UserProfileStorage struct {
	db *pgxpool.Pool
}

// NewUserProfileStorage creates postgres storage
func NewUserProfileStorage(db *pgxpool.Pool) *UserProfileStorage {
	return &UserProfileStorage{
		db: db,
	}
}

func (up *UserProfileStorage) GetUserProfileByID(ctx context.Context, userID int) (entities.UserProfile, error) {
	var user []entities.UserProfile

	err := pgxscan.Select(ctx, up.db, &user, `SELECT user_id, username, bio, avatar FROM profile_data WHERE user_id = $1`, userID)
	if err != nil {
		return entities.UserProfile{}, err
	}

	if user == nil {
		return entities.UserProfile{}, errors.New("User not found")
	}
	return user[0], nil
}

func (up *UserProfileStorage) EditUsername(ctx context.Context, userID int, username string) error {
	_, err := up.db.Exec(ctx, `UPDATE profile_data SET username = $1 WHERE user_id = $2`, username, userID)
	return err
}

func (up *UserProfileStorage) EditUserBio(ctx context.Context, userID int, bio string) error {
	_, err := up.db.Exec(ctx, `UPDATE profile_data SET bio = $1 WHERE user_id = $2`, bio, userID)
	return err
}

func (up *UserProfileStorage) EditUserAvatar(ctx context.Context, userID int, avatar string) error {
	_, err := up.db.Exec(ctx, "UPDATE profile_data SET avatar = $1 WHERE user_id = $2", avatar, userID)
	return err
}
