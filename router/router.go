package router

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/logout"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/login"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/cors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

func SetupRouter(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(cors.CorsMiddleware)

	router.Mount("/sights", SightRoutes())
	router.Mount("/signup", SignUpRoutes())
	router.Mount("/login", AuthRoutes())
	router.Mount("/logout", LogOutRoutes())

	return router
}

func SightRoutes() chi.Router {
	router := chi.NewRouter()
	sightsHandler := sight.SightsHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Sight, sight.Sights]{ServeHTTP: sightsHandler.GetSights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router

}

func SignUpRoutes() chi.Router {
	router := chi.NewRouter()

	regHandler := registration.Registration{}
	wrapperInstance := &wrapper.Wrapper[entities.User, registration.Response]{ServeHTTP: regHandler.SignUp}
	router.Post("/", wrapperInstance.HandlerWrapper)

	router.Mount("/sights", SightRoutes())

	return router
}

func LogOutRoutes() chi.Router {
	router := chi.NewRouter()

	logOutHandler := logout.Logout{}
	wrapperInstance := &wrapper.Wrapper[entities.User, logout.Response]{ServeHTTP: logOutHandler.LogOut}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func AuthRoutes() chi.Router {
	router := chi.NewRouter()

	authHandler := login.Authorization{}
	wrapperInstance := &wrapper.Wrapper[entities.User, login.Response]{ServeHTTP: authHandler.Authorize}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}
