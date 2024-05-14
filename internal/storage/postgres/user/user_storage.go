package user

import (
	"context"
	"sync"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserStorage struct
type UserStorage struct {
	db  *pgxpool.Pool
	mu  sync.Mutex
	ctx context.Context
}

// NewUserStorage creates postgres storage
func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	ctx, cancel := context.WithCancel(context.Background())
	// Обеспечение освобождения ресурсов контекста при завершении работы
	go func() {
		<-ctx.Done()
		cancel()
	}()

	return &UserStorage{
		db:  db,
		ctx: ctx,
		mu:  sync.Mutex{},
	}
}

func (us *UserStorage) SaveUser(email, password, salt string) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.ctx = context.Background()
	_, err := us.db.Exec(us.ctx, `INSERT INTO user_data (email, passwrd, salt) VALUES ($1, $2, $3)`, email, password, salt)
	return err
}

func (us *UserStorage) ChangeEmail(userID int, email string) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.ctx = context.Background()
	_, err := us.db.Exec(us.ctx, `UPDATE user_data SET email = $1 WHERE id = $2`, email, userID)
	return err
}

func (us *UserStorage) ChangePassword(userID int, password, salt string) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.ctx = context.Background()
	_, err := us.db.Exec(us.ctx, `UPDATE user_data SET passwrd = $1, salt = $2 WHERE id = $3`, password, salt, userID)
	return err
}

func (us *UserStorage) GetUserByID(id int) (entities.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.ctx = context.Background()
	var user entities.User
	err := us.db.QueryRow(us.ctx, `SELECT id, email, passwrd, salt FROM user_data WHERE id = $1`, id).Scan(&user.ID, &user.Username, &user.Password, &user.Salt)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (us *UserStorage) GetUserByEmail(email string) (entities.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.ctx = context.Background()
	var user entities.User
	err := us.db.QueryRow(us.ctx, `SELECT id, email, passwrd, salt FROM user_data WHERE email = $1`, email).Scan(&user.ID, &user.Username, &user.Password, &user.Salt)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (us *UserStorage) DeleteUser(userID int) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.ctx = context.Background()
	_, err := us.db.Exec(us.ctx, `DELETE FROM user_data WHERE id = $1`, userID)
	return err
}

func (us *UserStorage) IsEmailTaken(email string) (bool, error) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.ctx = context.Background()

	var count int
	err := us.db.QueryRow(us.ctx, `SELECT COUNT(*) FROM user_data WHERE email = $1`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
