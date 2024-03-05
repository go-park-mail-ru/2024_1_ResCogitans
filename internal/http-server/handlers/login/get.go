package login

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/models/user"
	_ "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/response"
	"github.com/gorilla/securecookie"
	"net/http"
)

type Authorization struct{} //Обработчик запросов на авторизацию

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func (h *Authorization) ServeHTTP(ctx context.Context, _ user.User) (response.Response, error) {
	username := ctx.Value("username").(string)
	password := ctx.Value("password").(string)

	if user.UserValidation(username, password) {
		setSession(username, w)

		return response.OK(nil), nil
	}

	return response.GetError("Authentication failed", nil), nil
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
