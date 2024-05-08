package httputils

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

const (
	ResponseWriterKey     = "responseWriter"
	HttpRequestKey        = "HttpRequest"
	RequestPathParamsKey  = "requestPathParams"
	RequestQueryParamsKey = "requestQueryParams"
)

func SetRequestToCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, HttpRequestKey, r)
}

func GetRequestFromCtx(ctx context.Context) (*http.Request, error) {
	r, ok := ctx.Value(HttpRequestKey).(*http.Request)
	if !ok {
		return nil, errors.New("Failed getting request")
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

func GetResponseWriterFromCtx(ctx context.Context) (http.ResponseWriter, error) {
	w, ok := ctx.Value(ResponseWriterKey).(http.ResponseWriter)
	if !ok {
		return nil, errors.New("Failed getting response writer")
	}
	return w, nil
}

func GetUserFromCtx(ctx context.Context) (int, error) {
	user := ctx.Value("userID")
	userID, ok := user.(int)
	if !ok {
		return 0, errors.New("Failed getting user from context")
	}
	return userID, nil
}

func SetQueryParamToCtx(ctx context.Context, queryParams map[string]string) context.Context {
	return context.WithValue(ctx, RequestQueryParamsKey, queryParams)
}

func GetQueryParamsFromCtx(ctx context.Context) map[string]string {
	queryParams, ok := ctx.Value(RequestQueryParamsKey).(map[string]string)
	if !ok {
		return nil
	}
	return queryParams
}
