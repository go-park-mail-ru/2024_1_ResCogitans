package authorization

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"
)

type AuthorizationHandler struct{}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Response struct {
	User      UserResponse
	SessionID string
}

var (
	CookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
	SessionStore = make(map[string]string)
	mu           sync.Mutex
)

func (h *AuthorizationHandler) Authorize(ctx context.Context, requestData entities.User) (Response, error) {
	username := requestData.Username
	password := requestData.Password

	responseWriter, ok := httputils.ContextWriter(ctx)
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
		if err := SetSession(sessionID, responseWriter); err != nil {
			return Response{}, errors.Wrap(err, "failed to set session cookie")
		}

		userResponse := UserResponse{
			ID:       user.ID,
			Username: user.Username,
		}

		return Response{User: userResponse, SessionID: sessionID}, nil
	}

	return Response{}, errors.New("Authorization failed")
}

func SetSession(sessionID string, response http.ResponseWriter) error {
	encoded, err := CookieHandler.Encode("session", sessionID)
	if err != nil {
		return errors.Wrap(err, "failed to encode session cookie")
	}

	cookie := &http.Cookie{
		Name:    "session",
		Value:   encoded,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(response, cookie)
	return nil
}

func RevokeSession(sessionID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(SessionStore, sessionID)
}

func (h *AuthorizationHandler) LogOut(ctx context.Context, requestData entities.User) (Response, error) {
	var r *http.Request
	if val, ok := ctx.Value(httputils.HttpRequestKey).(http.Request); ok {
		r = &val
	} else {
		return Response{}, errors.New("failed getting request")
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		return Response{}, errors.Wrap(err, "failed getting session cookie")
	}

	var sessionID string
	if err := CookieHandler.Decode("session", cookie.Value, &sessionID); err != nil {
		return Response{}, errors.Wrap(err, "failed decoding session cookie")
	}

	mu.Lock()
	_, sessionExists := SessionStore[sessionID]
	mu.Unlock()

	if !sessionExists {
		return Response{}, errors.New("User is not authorized or session has already been revoked")
	}

	RevokeSession(sessionID)

	w, ok := httputils.ContextWriter(ctx)
	if !ok {
		return Response{}, errors.New("Internal Error")
	}
	cookie = &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return Response{}, nil
}
