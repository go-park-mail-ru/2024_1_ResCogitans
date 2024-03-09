package httputils

import (
	"context"
	"net/http"
)

const (
	ResponseWriterKey    = "responseWriter"
	HttpRequestKey       = "HttpRequest"
	RequestPathParamsKey = "requestPathParams"
)

func ContextWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(ResponseWriterKey).(http.ResponseWriter)
	return w, ok
}

func HttpRequest(ctx context.Context) (*http.Request, bool) {
	w, ok := ctx.Value(HttpRequestKey).(*http.Request)
	return w, ok
}
