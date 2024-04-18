package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
)

type UseCases struct {
	UserUseCase    usecase.UserUseCaseInterface
	SessionUseCase usecase.SessionInterface
	ProfileUseCase usecase.ProfileUseCaseInterface
	SightUseCase   usecase.SightUseCaseInterface
	JourneyUseCase usecase.JourneyUseCaseInterface
	CommentUseCase usecase.CommentUseCaseInterface
}

func UseCaseInit(storages *Storages) *UseCases {
	return &UseCases{
		UserUseCase:    usecase.NewUserUseCase(storages.UserStorage),
		SessionUseCase: usecase.NewSessionUseCase(storages.SessionStorage),
		ProfileUseCase: usecase.NewProfileUseCase(storages.ProfileStorage),
		SightUseCase:   usecase.NewSightUseCase(storages.SightStorage),
		JourneyUseCase: usecase.NewJourneyUseCase(storages.SightStorage),
		CommentUseCase: usecase.NewCommentUseCase(storages.SightStorage),
	}
}
