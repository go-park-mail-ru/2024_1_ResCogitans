package httperrors

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrBadRequest = HttpError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid request data"),
	}

	ErrUnauthorized = HttpError{
		Code:    http.StatusUnauthorized,
		Message: errors.New("authorization error"),
	}

	ErrForbidden = HttpError{
		Code:    http.StatusForbidden,
		Message: errors.New("access error"),
	}

	ErrNotFound = HttpError{
		Code:    http.StatusNotFound,
		Message: errors.New("page not found"),
	}

	ErrInternal = HttpError{
		Code:    http.StatusInternalServerError,
		Message: errors.New("internal server error"),
	}
)
