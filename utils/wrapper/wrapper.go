package wrapper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type Handler[T Validator, R Response] interface {
	ServeHTTP(context.Context, T) (R, error)
}

type Validator interface {
	Validate() error
}

type Response interface {
	GetStatus() string
	GetMessage() string
	GetData() interface{}
}

func HandlerWrapper[T Validator, R Response](handler Handler[T, R]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		pathParams := GetPathParams(r)
		ctx = context.WithValue(ctx, "requestPathParams", pathParams)

		limitedReader := io.LimitReader(r.Body, 1_000_000)

		var requestData T
		err := json.NewDecoder(limitedReader).Decode(&requestData)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		if err := requestData.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, "requestData", requestData)

		response, err := handler.ServeHTTP(ctx, requestData)
		if err != nil {
			logger.Logger().Error("Handler error", "error", err)

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
	}
}

func GetPathParams(r *http.Request) map[string]string {
	params := chi.RouteContext(r.Context()).URLParams
	pathParams := make(map[string]string)
	for k := len(params.Keys) - 1; k >= 0; k-- {
		key := params.Keys[k]
		value := params.Values[k]
		pathParams[key] = value
	}
	return pathParams
}
