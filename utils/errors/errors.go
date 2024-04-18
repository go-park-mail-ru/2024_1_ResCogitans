package httperrors

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/pkg/errors"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e HttpError) Error() string {
	return e.Message
}

func NewHttpError(code int, message string) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

func UnwrapHttpError(err error) HttpError {
	var httpError HttpError
	if errors.As(err, &httpError) {
		return NewHttpError(httpError.Code, httpError.Message)
	}
	return NewHttpError(http.StatusInternalServerError, err.Error())
}

func IsHttpError(err error) bool {
	var httpError HttpError
	if errors.As(err, &httpError) {
		return true
	}
	return false
}

var errInternalBytes = []byte(`{"error": "internal error"}`)

func WriteHttpError(errIn error, w http.ResponseWriter) {
	httpError := UnwrapHttpError(errIn)
	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(httpError)
	if err != nil {
		logger.Logger().Error("error marshal err", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write(errInternalBytes); writeErr != nil {
			logger.Logger().Error("error writing fallback error", "error", writeErr.Error())
		}
		return
	}

	w.WriteHeader(httpError.Code)
	if _, writeErr := w.Write(bytes); writeErr != nil {
		logger.Logger().Error("error writing http error", "error", writeErr)
		return
	}
}
