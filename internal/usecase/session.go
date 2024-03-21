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

var SessionStore = make(map[string]int)
var mu sync.Mutex

func SetSession(w http.ResponseWriter, userID int) error {
	sessionID := uuid.New().String()
	mu.Lock()
	SessionStore[sessionID] = userID
	mu.Unlock()
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

func GetSession(r *http.Request) int {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return 0
	}

	var sessionID string
	if err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID); err == nil {
		mu.Lock()
		defer mu.Unlock()
		return SessionStore[sessionID]
	}
	return 0
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return err
	}

	var sessionID string
	err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID)
	if err != nil {
		return err
	}

	mu.Lock()
	delete(SessionStore, sessionID)
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   sessionId,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return nil
}
