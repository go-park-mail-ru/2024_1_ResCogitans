package wrapper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request) (interface{}, error)
	Validate() error
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func HandlerWrapper(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		params := chi.RouteContext(r.Context()).URLParams
		for k := len(params.Keys) - 1; k >= 0; k-- {
			key := params.Keys[k]
			value := params.Values[k]
			ctx = context.WithValue(ctx, key, value)
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1000000)

		if err := handler.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := handler.ServeHTTP(ctx, w, r)
		if err != nil {
			logger.Logger().Error("Handler error", "error", err)

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
	}
}
