package httputils

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

const (
	ResponseWriterKey    = "responseWriter"
	HttpRequestKey       = "HttpRequest"
	RequestPathParamsKey = "requestPathParams"
)

func SetRequestToCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, HttpRequestKey, r)
}

func GetRequestFromCtx(ctx context.Context) (*http.Request, error) {
	r, ok := ctx.Value(HttpRequestKey).(*http.Request)
	if !ok {
		return nil, errors.New("failed getting request")
	}
	return r, nil
}

func SetPathParamsToCtx(ctx context.Context, pathParams map[string]string) context.Context {
	return context.WithValue(ctx, RequestPathParamsKey, pathParams)
}

func GetPathParamsFromCtx(ctx context.Context) map[string]string {
	pathParams, ok := ctx.Value(RequestPathParamsKey).(map[string]string)
	if !ok {
		return nil
	}
	return pathParams
}

func SetResponseWriterToCtx(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, ResponseWriterKey, w)
}

// GetResponseWriterFromCtx извлекает http.ResponseWriter из контекста.
func GetResponseWriterFromCtx(ctx context.Context) (http.ResponseWriter, error) {
	w, ok := ctx.Value(ResponseWriterKey).(http.ResponseWriter)
	if !ok {
		return nil, errors.New("failed getting response writer")
	}
	return w, nil
}
