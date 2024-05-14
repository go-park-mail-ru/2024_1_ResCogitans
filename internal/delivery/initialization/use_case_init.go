package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
)

type UseCases struct {
	UserUseCase     *usecase.UserUseCase
	SessionUseCase  *usecase.SessionUseCase
	ProfileUseCase  *usecase.ProfileUseCase
	SightUseCase    *usecase.SightUseCase
	JourneyUseCase  *usecase.JourneyUseCase
	CommentUseCase  *usecase.CommentUseCase
	QuestionUseCase *usecase.QuestionUseCase
	CSRFUseCase     *usecase.CSRFUseCase
}

func UseCaseInit(storages *Storages) *UseCases {
	return &UseCases{
		UserUseCase:     usecase.NewUserUseCase(storages.UserStorage),
		ProfileUseCase:  usecase.NewProfileUseCase(storages.ProfileStorage),
		SightUseCase:    usecase.NewSightUseCase(storages.SightStorage),
		JourneyUseCase:  usecase.NewJourneyUseCase(storages.SightStorage),
		CommentUseCase:  usecase.NewCommentUseCase(storages.SightStorage),
		QuestionUseCase: usecase.NewQuestionUseCase(storages.QuestionStorage),
		SessionUseCase:  usecase.NewSessionUseCase(storages.SessionStorage),
		CSRFUseCase:     usecase.NewCSRFUseCase(storages.CSRFStorage),
	}
}
