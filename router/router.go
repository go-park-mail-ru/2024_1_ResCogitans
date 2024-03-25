package router

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/updateUserData"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/app"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/cors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

func SetupRouter(app *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(cors.CorsMiddleware)
	router.Use(app.AuthMiddleware.Auth)

	router.Mount("/sights", SightRoute())
	router.Mount("/signup", SignUpRoute(app.RegHandler))
	router.Mount("/login", AuthRoute(app.AuthHandler))
	router.Mount("/logout", LogOutRoute(app.AuthHandler))
	router.Mount("/updatedata", UpdateDateRoute(app.UpdateUserDataHandler))

	return router
}

func SightRoute() chi.Router {
	router := chi.NewRouter()
	sightsHandler := sight.SightsHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Sight, sight.Sights]{ServeHTTP: sightsHandler.GetSights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func CreateRoute[T wrapper.Validator, Resp any](f wrapper.ServeHTTPFunc[T, Resp]) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[T, Resp]{ServeHTTP: f}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func CreateServeHTTPFunc[T wrapper.Validator, R any](f func(ctx context.Context, request T) (R, error)) wrapper.ServeHTTPFunc[T, R] {
	return f
}

func SignUpRoute(regHandler *registration.RegistrationHandler) chi.Router {
	ServeHTTPFunc := CreateServeHTTPFunc[entities.User, entities.UserResponse](regHandler.SignUp)
	return CreateRoute[entities.User, entities.UserResponse](ServeHTTPFunc)
}

func LogOutRoute(authHandler *authorization.AuthorizationHandler) chi.Router {
	ServeHTTPFunc := CreateServeHTTPFunc[entities.User, entities.UserResponse](authHandler.LogOut)
	return CreateRoute[entities.User, entities.UserResponse](ServeHTTPFunc)
}

func AuthRoute(authHandler *authorization.AuthorizationHandler) chi.Router {
	ServeHTTPFunc := CreateServeHTTPFunc[entities.User, entities.UserResponse](authHandler.Authorize)
	return CreateRoute[entities.User, entities.UserResponse](ServeHTTPFunc)
}

func UpdateDateRoute(updateHandler *updateUserData.UpdateDataHandler) chi.Router {
	ServeHTTPFunc := CreateServeHTTPFunc[entities.User, entities.UserResponse](updateHandler.Update)
	return CreateRoute[entities.User, entities.UserResponse](ServeHTTPFunc)
}
