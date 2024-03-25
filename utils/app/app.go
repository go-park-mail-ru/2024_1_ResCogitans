package app

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
)

type App struct {
	AuthHandler    *authorization.AuthorizationHandler
	RegHandler     *registration.RegistrationHandler
	AuthMiddleware *middle.AuthMiddleware
}

func NewApp(
	authHandler *authorization.AuthorizationHandler,
	regHandler *registration.RegistrationHandler,
	authMiddleware *middle.AuthMiddleware,
) *App {
	return &App{
		AuthHandler:    authHandler,
		RegHandler:     regHandler,
		AuthMiddleware: authMiddleware,
	}
}
