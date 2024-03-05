package login

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/pkg/errors"

	_ "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/gorilla/securecookie"
	"net/http"
)

type Authorization struct{}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type contextKey string

const responseWriterKey = contextKey("responseWriter")

func ContextWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	return w, ok
}

func (h *Authorization) Authorize(ctx context.Context, _ entities.User) (*entities.User, error) {
	username := ctx.Value("username").(string)
	password := ctx.Value("password").(string)
	responseWriter, ok := ContextWriter(ctx)
	if !ok {
		return nil, errors.New("Response Writer not found in context")
	}

	if entities.UserValidation(username, password) {
		setSession(username, responseWriter)

		authorizedUser, err := entities.GetUserByUsername(username)
		if err != nil {
			return nil, errors.New("Problem with searching for a profile by username")
		}

		return authorizedUser, nil
	}

	return nil, errors.New("Authorization failed")
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}
