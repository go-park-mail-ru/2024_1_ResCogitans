package usecase

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

const sessionId = "session_id"

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type AuthUseCase struct {
	SessionStore map[string]int
	mu           sync.Mutex
}

var Auth *AuthUseCase

func init() {
	Auth = &AuthUseCase{
		SessionStore: make(map[string]int),
	}
}

func (a *AuthUseCase) SetSession(w http.ResponseWriter, userID int) error {
	sessionID := uuid.New().String()
	a.mu.Lock()
	a.SessionStore[sessionID] = userID
	a.mu.Unlock()
	encoded, err := CookieHandler.Encode(sessionId, sessionID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    sessionId,
		Value:   encoded,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
	return nil
}

func (a *AuthUseCase) GetSession(r *http.Request) int {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return 0
	}

	var sessionID string
	if err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID); err == nil {
		a.mu.Lock()
		defer a.mu.Unlock()
		return a.SessionStore[sessionID]
	}
	return 0
}

func (a *AuthUseCase) ClearSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return err
	}

	var sessionID string
	err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID)
	if err != nil {
		return err
	}

	a.mu.Lock()
	delete(a.SessionStore, sessionID)
	a.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   sessionId,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return nil
}
