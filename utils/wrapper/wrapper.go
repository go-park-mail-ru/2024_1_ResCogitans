package wrapper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

const (
	requestPathParamsKey = "requestPathParams"
	requestDataKey       = "requestData"
)

var (
	validationErr = errors.HttpError{
		Code:    http.StatusBadRequest,
		Message: "invalid request data",
	}

	decodingErr = errors.HttpError{
		Code:    http.StatusBadRequest,
		Message: "json decoding error",
	}

	encodingErr = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "json encoding error",
	}
)

type Handler[T Validator, Resp any] interface {
	ServeHTTP(context.Context, T) (Resp, error)
}

type Validator interface {
	Validate() error
}

func HandlerWrapper[T Validator, Resp any](handler func(ctx context.Context, req T) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.Logger()

		pathParams := GetPathParams(r)
		ctx = SetPathParamsToCtx(ctx, pathParams)

		limitedReader := io.LimitReader(r.Body, 1_000_000)

		var requestData T
		err := json.NewDecoder(limitedReader).Decode(&requestData)
		// if err != nil {
		// 	logger.Error("Error decoding request body", "error", err)
		// 	errors.WriteHttpError(decodingErr, w)
		// 	return
		// }

		if err := requestData.Validate(); err != nil {
			logger.Error("Validation error", "error", err)
			errors.WriteHttpError(validationErr, w)
			return
		}

		ctx = context.WithValue(ctx, requestDataKey, requestData)

		response, err := handler(ctx, requestData)
		if err != nil {
			logger.Error("Handler error", "error", err)
			errors.WriteHttpError(errors.HttpError{Code: http.StatusInternalServerError, Message: err.Error()}, w)
			return
		}

		rawJSON, err := json.Marshal(response)
		if err != nil {
			logger.Error("Error encoding response", "error", err)
			errors.WriteHttpError(encodingErr, w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(rawJSON)
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

func SetPathParamsToCtx(ctx context.Context, pathParams map[string]string) context.Context {
	return context.WithValue(ctx, requestPathParamsKey, pathParams)
}

func GetPathParamsFromCtx(ctx context.Context) map[string]string {
	pathParams, ok := ctx.Value(requestPathParamsKey).(map[string]string)
	if !ok {
		return nil
	}
	return pathParams
}
