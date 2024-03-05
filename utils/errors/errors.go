package errors

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/pkg/errors"
)

type HttpError struct {
	Code    int
	Message string `json:"error"`
}

var (
	errFallBack = HttpError{
		Code:    http.StatusInternalServerError,
		Message: "internal error",
	}
	errInternalBytes = []byte("{\n\t\"error\": \"internal error\"\n}")
)

func (e HttpError) Error() string {
	return e.Message
}

func WriteHttpError(errIn error, w http.ResponseWriter) {
	var httpErr HttpError
	if !errors.As(errIn, &httpErr) {
		httpErr = errFallBack
	}

	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(httpErr)
	if err != nil {
		logger.Logger().Error("error marshal err", "error", errIn)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errInternalBytes)
		return
	}

	w.WriteHeader(httpErr.Code)
	w.Write(bytes)
}
