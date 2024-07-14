package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/comment"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/journey"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/question"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/user"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/redis/csrf"
)

type Storages struct {
	UserStorage     *user.UserStorage
	ProfileStorage  *user.UserProfileStorage
	SightStorage    *sight.SightStorage
	CommentStorage  *comment.CommentStorage
	JourneyStorage  *journey.JourneyStorage
	QuestionStorage *question.QuestionStorage
	CSRFStorage     *csrf.CSRFStorage
}

func StorageInit(DBs *DBs) *Storages {
	return &Storages{
		UserStorage:     user.NewUserStorage(DBs.PostgresDB),
		ProfileStorage:  user.NewUserProfileStorage(DBs.PostgresDB),
		SightStorage:    sight.NewSightStorage(DBs.PostgresDB),
		CommentStorage:  comment.NewCommentStorage(DBs.PostgresDB),
		JourneyStorage:  journey.NewJourneyStorage(DBs.PostgresDB),
		QuestionStorage: question.NewQuestionStorage(DBs.PostgresDB),
		CSRFStorage:     csrf.NewCSRFStorage(DBs.CsrfDB),
	}
}
