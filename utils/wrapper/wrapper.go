package wrapper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type Validator interface {
	Validate() error
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type contextKey string

func HandlerWrapper(handler func(context.Context, http.ResponseWriter, *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// TODO: parse path params

		// const queryParamsKey contextKey = "id"
		// queryParams := r.URL.Query()
		// logger.Logger().Info("queryParams:", queryParams)
		// ctx = context.WithValue(ctx, queryParamsKey, queryParams)

		r.Body = http.MaxBytesReader(w, r.Body, 1000000)

		var requestBody Validator
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := requestBody.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := handler(ctx, w, r)
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
