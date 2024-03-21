package errors

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

var (
	errFallBack = HttpError{
		Code:    http.StatusInternalServerError,
		Message: "internal error",
	}
	errInternalBytes = []byte(`{"error": "internal error"}`)
)

func (e HttpError) Error() string {
	return e.Message
}

func WriteHttpError(errIn error, w http.ResponseWriter) error {
	var httpErr HttpError
	if !errors.As(errIn, &httpErr) {
		httpErr = errFallBack
	}

	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(httpErr)
	if err != nil {
		logger.Logger().Error("error marshal err", "error", errIn)
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write(errInternalBytes); writeErr != nil {
			logger.Logger().Error("error writing fallback error", "error", writeErr)
		}
		return err
	}

	w.WriteHeader(httpErr.Code)
	if _, writeErr := w.Write(bytes); writeErr != nil {
		logger.Logger().Error("error writing http error", "error", writeErr)
		return writeErr
	}
	return nil
}
