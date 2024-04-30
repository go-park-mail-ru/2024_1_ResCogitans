package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/category"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/question"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/user"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/redis/session"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storages struct {
	UserStorage     storage.UserStorageInterface
	ProfileStorage  storage.UserProfileStorageInterface
	SessionStorage  storage.SessionStorageInterface
	SightStorage    storage.SightStorageInterface
	QuestionStorage storage.QuestionInterface
	CategoryStorage storage.CategoryStorageInterface
}

func StorageInit(pdb *pgxpool.Pool, rdb *redis.Client) *Storages {
	return &Storages{
		UserStorage:     user.NewUserStorage(pdb),
		ProfileStorage:  user.NewUserProfileStorage(pdb),
		SessionStorage:  session.NewSessionStorage(rdb),
		SightStorage:    sight.NewSightStorage(pdb),
		QuestionStorage: question.NewQuestionStorage(pdb),
		CategoryStorage: category.NewCategoryStorage(pdb),
	}
}
