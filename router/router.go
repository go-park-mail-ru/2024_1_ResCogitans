package router

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/sight"
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

	router.Mount("/sights", SightRoutes())
	router.Mount("/signup", SignUpRoutes(app.RegHandler))
	router.Mount("/login", AuthRoutes(app.AuthHandler))
	router.Mount("/logout", LogOutRoutes(app.AuthHandler))

	return router
}

func SightRoutes() chi.Router {
	router := chi.NewRouter()
	sightsHandler := sight.SightsHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Sight, sight.Sights]{ServeHTTP: sightsHandler.GetSights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func SignUpRoutes(regHandler *registration.RegistrationHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.User, registration.UserResponse]{ServeHTTP: regHandler.SignUp}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func LogOutRoutes(authHandler *authorization.AuthorizationHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.User, authorization.UserResponse]{ServeHTTP: authHandler.LogOut}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func AuthRoutes(authHandler *authorization.AuthorizationHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.User, authorization.UserResponse]{ServeHTTP: authHandler.Authorize}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}
