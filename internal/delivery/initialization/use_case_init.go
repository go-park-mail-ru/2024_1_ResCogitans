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
	AlbumUseCase    usecase.AlbumUseCaseInterface
}

func UseCaseInit(storages *Storages, conn *grpc.ClientConn) *UseCases {
	return &UseCases{
		UserUseCase:     usecase.NewUserUseCase(storages.UserStorage),
		SessionUseCase:  usecase.NewSessionUseCase(conn),
		ProfileUseCase:  usecase.NewProfileUseCase(storages.ProfileStorage),
		SightUseCase:    usecase.NewSightUseCase(storages.SightStorage),
		JourneyUseCase:  usecase.NewJourneyUseCase(storages.SightStorage),
		CommentUseCase:  usecase.NewCommentUseCase(storages.SightStorage),
		QuestionUseCase: usecase.NewQuestionUseCase(storages.QuestionStorage),
		AlbumUseCase:    usecase.NewAlbumUseCase(storages.AlbumStorage),
	}
}
