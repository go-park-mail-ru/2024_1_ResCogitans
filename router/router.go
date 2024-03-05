package router

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

type RequestData struct {
}

func (rd RequestData) Validate() error {
	return nil
}

func myHandler(ctx context.Context) (string, error) {
	return "Hello World", nil
}

func SetupRouter(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	wrapperInstance := &wrapper.Wrapper[RequestData, string]{ServeHTTP: myHandler}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}
