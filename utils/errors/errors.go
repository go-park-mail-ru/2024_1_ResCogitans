package httperrors

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type HttpError struct {
	Code    int
	Message error
}

type HttpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func NewHttpError(code int, message error) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

var errInternalBytes = []byte(`{"error": "internal error"}`)

func (e HttpError) Error() error {
	return e.Message
}

func WriteHttpError(errIn HttpError, w http.ResponseWriter) {
	var response HttpResponse
	response.Code = errIn.Code
	response.Message = errIn.Message.Error()
	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(response)
	if err != nil {
		logger.Logger().Error("error marshal err", "error", response.Message)
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write(errInternalBytes); writeErr != nil {
			logger.Logger().Error("error writing fallback error", "error", writeErr)
		}
		return
	}

	w.WriteHeader(response.Code)
	if _, writeErr := w.Write(bytes); writeErr != nil {
		logger.Logger().Error("error writing http error", "error", writeErr)
		return
	}
}
