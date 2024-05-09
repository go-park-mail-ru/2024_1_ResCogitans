package httputils

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type contextKey string

const (
	keyResponseWriter     contextKey = "responseWriter"
	keyHttpRequest        contextKey = "HttpRequest"
	keyRequestPathParams  contextKey = "requestPathParams"
	keyRequestQueryParams contextKey = "requestQueryParams"
	keyUserID             contextKey = "userID"
)

func SetRequestToCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, keyHttpRequest, r)
}

func GetRequestFromCtx(ctx context.Context) (*http.Request, error) {
	r, ok := ctx.Value(keyHttpRequest).(*http.Request)
	if !ok {
		return nil, errors.New("failed getting request")
	}
	return r, nil
}

func SetPathParamsToCtx(ctx context.Context, pathParams map[string]string) context.Context {
	return context.WithValue(ctx, keyRequestPathParams, pathParams)
}

func GetPathParamsFromCtx(ctx context.Context) map[string]string {
	pathParams, ok := ctx.Value(keyRequestPathParams).(map[string]string)
	if !ok {
		return nil
	}
	return pathParams
}

func SetResponseWriterToCtx(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, keyResponseWriter, w)
}

func GetResponseWriterFromCtx(ctx context.Context) (http.ResponseWriter, error) {
	w, ok := ctx.Value(keyResponseWriter).(http.ResponseWriter)
	if !ok {
		return nil, errors.New("failed getting response writer")
	}
	return w, nil
}

func GetUserFromCtx(ctx context.Context) (int, error) {
	user := ctx.Value(keyUserID)
	userID, ok := user.(int)
	if !ok {
		return 0, errors.New("failed getting user from context")
	}
	return userID, nil
}

func SetQueryParamToCtx(ctx context.Context, queryParams map[string]string) context.Context {
	return context.WithValue(ctx, keyRequestQueryParams, queryParams)
}

func GetQueryParamsFromCtx(ctx context.Context) map[string]string {
	queryParams, ok := ctx.Value(keyRequestQueryParams).(map[string]string)
	if !ok {
		return nil
	}
	return queryParams
}
