package router

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/sight"
	modelSight "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/models/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

func SetupRouter(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	router.Mount("/sights", SightRoutes())

	return router
}

func SightRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", wrapper.HandlerWrapper[modelSight.Sight](&sight.GetSights{}))

	return r
}
