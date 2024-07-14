package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
)

type UseCases struct {
	UserUseCase     usecase.UserUseCaseInterface
	ProfileUseCase  usecase.ProfileUseCaseInterface
	SightUseCase    usecase.SightUseCaseInterface
	JourneyUseCase  usecase.JourneyUseCaseInterface
	CommentUseCase  usecase.CommentUseCaseInterface
	QuestionUseCase usecase.QuestionUseCaseInterface
	CSRFUseCase     usecase.CSRFInterface
	SessionUseCase  usecase.SessionInterface
}

func UseCaseInit(storages *Storages, sessionUseCase *usecase.SessionUseCase) *UseCases {
	return &UseCases{
		UserUseCase:     usecase.NewUserUseCase(storages.UserStorage),
		ProfileUseCase:  usecase.NewProfileUseCase(storages.ProfileStorage),
		SightUseCase:    usecase.NewSightUseCase(storages.SightStorage, storages.CommentStorage),
		JourneyUseCase:  usecase.NewJourneyUseCase(storages.JourneyStorage),
		CommentUseCase:  usecase.NewCommentUseCase(storages.CommentStorage),
		QuestionUseCase: usecase.NewQuestionUseCase(storages.QuestionStorage),
		CSRFUseCase:     usecase.NewCSRFUseCase(storages.CSRFStorage),
		SessionUseCase:  sessionUseCase,
	}
}
