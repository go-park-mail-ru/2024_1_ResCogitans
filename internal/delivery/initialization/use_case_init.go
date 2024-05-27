package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"google.golang.org/grpc"
)

type UseCases struct {
	UserUseCase     usecase.UserUseCaseInterface
	SessionUseCase  usecase.SessionInterface
	ProfileUseCase  usecase.ProfileUseCaseInterface
	SightUseCase    usecase.SightUseCaseInterface
	JourneyUseCase  usecase.JourneyUseCaseInterface
	CommentUseCase  usecase.CommentUseCaseInterface
	QuestionUseCase usecase.QuestionUseCaseInterface
	CSRFUseCase     usecase.CSRFInterface
}

func UseCaseInit(storages *Storages, conn *grpc.ClientConn) *UseCases {
	return &UseCases{
		UserUseCase:     usecase.NewUserUseCase(storages.UserStorage),
		SessionUseCase:  usecase.NewSessionUseCase(conn),
		ProfileUseCase:  usecase.NewProfileUseCase(storages.ProfileStorage),
		SightUseCase:    usecase.NewSightUseCase(storages.SightStorage, storages.CommentStorage),
		JourneyUseCase:  usecase.NewJourneyUseCase(storages.JourneyStorage),
		CommentUseCase:  usecase.NewCommentUseCase(storages.CommentStorage),
		QuestionUseCase: usecase.NewQuestionUseCase(storages.QuestionStorage),
		CSRFUseCase:     usecase.NewCSRFUseCase(storages.CSRFStorage),
	}
}
