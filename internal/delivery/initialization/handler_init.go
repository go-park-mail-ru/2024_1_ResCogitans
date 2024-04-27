package initialization

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/comment"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/deactivation"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/journey"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/profile"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/quiz"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
)

type Handlers struct {
	AuthHandler         *authorization.AuthorizationHandler
	RegHandler          *registration.RegistrationHandler
	ProfileHandler      *profile.ProfileHandler
	DeactivationHandler *deactivation.DeactivationHandler
	SightHandler        *sight.SightHandler
	JourneyHandler      *journey.JourneyHandler
	CommentHandler      *comment.CommentHandler
	QuizHandler         *quiz.QuizHandler

	AuthMiddleware *middle.AuthMiddleware
}

func HandlerInit(cases *UseCases) *Handlers {
	return &Handlers{
		AuthHandler:         authorization.NewAuthorizationHandler(cases.SessionUseCase, cases.UserUseCase),
		RegHandler:          registration.NewRegistrationHandler(cases.SessionUseCase, cases.UserUseCase),
		ProfileHandler:      profile.NewProfileHandler(cases.ProfileUseCase),
		DeactivationHandler: deactivation.NewDeactivationHandler(cases.SessionUseCase, cases.UserUseCase),
		SightHandler:        sight.NewSightsHandler(cases.SightUseCase),
		JourneyHandler:      journey.NewJourneyHandler(cases.JourneyUseCase),
		CommentHandler:      comment.NewCommentHandler(cases.CommentUseCase),
		QuizHandler:         quiz.NewQuizHandler(cases.QuestionUseCase, cases.CommentUseCase, cases.JourneyUseCase),
		AuthMiddleware:      middle.NewAuthMiddleware(cases.SessionUseCase),
	}
}
