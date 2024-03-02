package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

type MyHandler struct{}

func (h *MyHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return "Hello World", nil
}

func (h *MyHandler) Validate() error {
	return nil
}

func SetupRouter(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	router.Get("/", wrapper.HandlerWrapper(&MyHandler{}))

	return router
}
