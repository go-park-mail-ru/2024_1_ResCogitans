package login

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"net/http"

	_ "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/gorilla/securecookie"
)

type Authorization struct{}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func ContextWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value("responseWriter").(http.ResponseWriter)
	return w, ok
}

func (h *Authorization) Authorize(ctx context.Context, _ entities.User) (Response, error) {
	requestData, ok := ctx.Value("requestData").(entities.User)
	if !ok {
		return Response{Status: http.StatusBadRequest, Message: "requestData not found in context"}, nil
	}

	username := requestData.Username
	password := requestData.Password

	responseWriter, ok := ContextWriter(ctx)
	if !ok {
		return Response{Status: http.StatusBadRequest, Message: "Response Writer not found in context"}, nil
	}

	println(username, "|", password)
	if entities.UserValidation(username, password) {

		_, err := entities.GetUserByUsername(username)
		if err != nil {
			return Response{Status: http.StatusBadRequest, Message: "Problem with searching for a profile by username"}, nil
		}

		setSession(username, responseWriter)

		return Response{Status: http.StatusOK}, nil
	}

	return Response{Status: http.StatusUnauthorized, Message: "Authorization failed"}, nil
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
