package entities

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var SessionStore = make(map[string]int)

func SetSession(w http.ResponseWriter, userID int) {
	sessionID := uuid.New().String()
	SessionStore[sessionID] = userID
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
}

func GetSession(r *http.Request) int {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0
	}
	userID, _ := SessionStore[cookie.Value]
	return userID
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		return
	}
	delete(SessionStore, cookie.Value)
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
