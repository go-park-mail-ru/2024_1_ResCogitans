package wrapper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type ServeHTTPFunc[T Validator, Resp any] func(ctx context.Context, request T) (Resp, error)

type Wrapper[T Validator, Resp any] struct {
	ServeHTTP ServeHTTPFunc[T, Resp]
}

type Validator interface {
	Validate() error
}

func (w *Wrapper[T, Resp]) HandlerWrapper(resWriter http.ResponseWriter, httpReq *http.Request) {
	ctx := httpReq.Context()
	logger := logger.Logger()

	pathParams := GetPathParams(httpReq)
	ctx = httputils.SetPathParamsToCtx(ctx, pathParams)
	ctx = httputils.SetResponseWriterToCtx(ctx, resWriter)
	ctx = httputils.SetRequestToCtx(ctx, httpReq)

	limitedReader := io.LimitReader(httpReq.Body, 1_000_000)

	var requestData T
	if httpReq.ContentLength > 0 {
		err := json.NewDecoder(limitedReader).Decode(&requestData)
		if err != nil {
			errText := "Error decoding request body"
			logger.Error(errText, "error", err)
			httperrors.WriteHttpError(err, resWriter)
			return
		}

		if err = requestData.Validate(); err != nil {
			errorText := "Validation error"
			logger.Error(errorText, "error", err)
			httperrors.WriteHttpError(err, resWriter)
			return
		}
	}

	response, err := w.ServeHTTP(ctx, requestData)
	if err != nil {
		logger.Error("Handler error", "error", err.Error())
		httperrors.WriteHttpError(err, resWriter)
		return
	}

	rawJSON, err := json.Marshal(response)
	if err != nil {
		logger.Error("Error encoding response", "error", err)
		httperrors.WriteHttpError(err, resWriter)
		return
	}

	resWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	resWriter.WriteHeader(http.StatusOK)
	_, _ = resWriter.Write(rawJSON)
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
