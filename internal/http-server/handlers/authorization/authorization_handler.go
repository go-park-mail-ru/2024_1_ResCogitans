package authorization

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"
)

type AuthorizationHandler struct{}

type Response struct {
	User      entities.User
	SessionID string
}

var (
	CookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
	SessionStore = make(map[string]string)
	mu           sync.Mutex
)

func ContextWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value("responseWriter").(http.ResponseWriter)
	return w, ok
}

func HttpRequest(ctx context.Context) (http.Request, bool) {
	w, ok := ctx.Value("HttpRequest").(http.Request)
	return w, ok
}

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (Response, error) {
	username := requestData.Username
	password := requestData.Password

	responseWriter, ok := ContextWriter(ctx)
	if !ok {
		return Response{}, errors.New("Internal Error")
	}

	if entities.UserValidation(username, password) {
		user, err := entities.GetUserByUsername(username)
		if err != nil {
			return Response{}, errors.Wrap(err, "Problem with searching for a profile by username")
		}

		sessionID := uuid.New().String()
		SessionStore[sessionID] = username
		SetSession(sessionID, responseWriter)

		return Response{User: *user, SessionID: sessionID}, nil
	}

	return Response{}, errors.New("Authorization failed")
}

func SetSession(sessionID string, response http.ResponseWriter) {
	if encoded, err := CookieHandler.Encode("session", sessionID); err == nil {
		cookie := &http.Cookie{
			Name:    "session",
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(response, cookie)
	}
}

func RevokeSession(sessionID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(SessionStore, sessionID)
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, requestData entities.User) (Response, error) {
	var r *http.Request
	logger.Logger().Info("SSADSD")
	if val, ok := ctx.Value("request").(http.Request); ok {
		r = &val
	} else {
		return Response{}, errors.New("failed getting request")
	}

	if cookie, err := r.Cookie("session"); err == nil {
		var sessionID string
		if err = CookieHandler.Decode("session", cookie.Value, &sessionID); err == nil {
			RevokeSession(sessionID)
		}
	}

	w, ok := ContextWriter(ctx)
	if !ok {
		return Response{}, errors.New("Internal Error")
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return Response{}, nil
}
