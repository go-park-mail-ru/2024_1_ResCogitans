package app

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/updateUserData"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
)

type App struct {
	AuthHandler           *authorization.AuthorizationHandler
	RegHandler            *registration.RegistrationHandler
	UpdateUserDataHandler *updateUserData.UpdateDataHandler
	AuthMiddleware        *middle.AuthMiddleware
}

func NewApp(
	authHandler *authorization.AuthorizationHandler,
	regHandler *registration.RegistrationHandler,
	updateHandler *updateUserData.UpdateDataHandler,
	authMiddleware *middle.AuthMiddleware,
) *App {
	return &App{
		AuthHandler:           authHandler,
		RegHandler:            regHandler,
		UpdateUserDataHandler: updateHandler,
		AuthMiddleware:        authMiddleware,
	}
}
