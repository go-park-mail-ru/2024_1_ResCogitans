package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/comment"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/journey"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/question"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/user"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/redis/csrf"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/redis/session"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storages struct {
	UserStorage     *user.UserStorage
	ProfileStorage  *user.UserProfileStorage
	SessionStorage  *session.RedisStorage
	SightStorage    *sight.SightStorage
	CommentStorage  *comment.CommentStorage
	JourneyStorage  *journey.JourneyStorage
	QuestionStorage *question.QuestionStorage
	CSRFStorage     *csrf.CSRFStorage
}

func StorageInit(pdb *pgxpool.Pool, rdb *redis.Client, cdb *redis.Client) *Storages {
	return &Storages{
		UserStorage:     user.NewUserStorage(pdb),
		ProfileStorage:  user.NewUserProfileStorage(pdb),
		SessionStorage:  session.NewSessionStorage(rdb),
		SightStorage:    sight.NewSightStorage(pdb),
		CommentStorage:  comment.NewCommentStorage(pdb),
		JourneyStorage:  journey.NewJourneyStorage(pdb),
		QuestionStorage: question.NewQuestionStorage(pdb),
		CSRFStorage:     csrf.NewCSRFStorage(cdb),
	}
}
