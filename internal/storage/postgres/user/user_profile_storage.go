package user

import (
	"context"
	"sync"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// UserProfileStorage struct
type UserProfileStorage struct {
	db  *pgxpool.Pool
	mu  sync.Mutex
	ctx context.Context
}

// NewUserProfileStorage creates postgres storage
func NewUserProfileStorage(db *pgxpool.Pool) storage.UserProfileStorageInterface {
	ctx, cancel := context.WithCancel(context.Background())
	// Обеспечение освобождения ресурсов контекста при завершении работы
	go func() {
		<-ctx.Done()
		cancel()
	}()

	return &UserProfileStorage{
		db:  db,
		ctx: ctx,
		mu:  sync.Mutex{},
	}
}

func (up *UserProfileStorage) GetUserProfileByID(userID int) (entities.UserProfile, error) {
	var user []entities.UserProfile
	up.mu.Lock()
	defer up.mu.Unlock()
	up.ctx = context.Background()

	err := pgxscan.Select(up.ctx, up.db, &user, `SELECT user_id, username, bio, avatar FROM "profile" WHERE user_id = $1`, userID)
	if err != nil {
		return entities.UserProfile{}, err
	}

	if user == nil {
		return entities.UserProfile{}, errors.New("User not found")
	}
	return user[0], nil
}

func (up *UserProfileStorage) EditUsername(userID int, username string) error {
	up.mu.Lock()
	defer up.mu.Unlock()
	up.ctx = context.Background()
	_, err := up.db.Exec(up.ctx, `UPDATE "profile" SET username = $1 WHERE user_id = $2`, username, userID)
	return err
}

func (up *UserProfileStorage) EditUserBio(userID int, bio string) error {
	up.mu.Lock()
	defer up.mu.Unlock()
	up.ctx = context.Background()
	_, err := up.db.Exec(up.ctx, `UPDATE "profile" SET bio = $1 WHERE user_id = $2`, bio, userID)
	return err
}

func (up *UserProfileStorage) EditUserAvatar(userID int, avatar string) error {
	up.mu.Lock()
	defer up.mu.Unlock()
	up.ctx = context.Background()
	_, err := up.db.Exec(up.ctx, `UPDATE "profile" SET avatar = $1 WHERE user_id = $2`, avatar, userID)
	return err
}
