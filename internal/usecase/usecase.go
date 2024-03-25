package usecase

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

const sessionId = "session_id"

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type AuthInterface interface {
	SetSession(w http.ResponseWriter, userID int) error
	GetSession(r *http.Request) (int, error)
	ClearSession(w http.ResponseWriter, r *http.Request) error
}

type AuthUseCase struct {
	SessionStorage storage.StorageInterface
}

func NewAuthUseCase(storage storage.StorageInterface) AuthInterface {
	return &AuthUseCase{
		SessionStorage: storage,
	}
}

func (a *AuthUseCase) SetSession(w http.ResponseWriter, userID int) error {
	sessionID := uuid.New().String()
	a.SessionStorage.SaveSession(sessionID, userID)
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

func (a *AuthUseCase) GetSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return 0, err
	}

	var sessionID string
	if err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID); err == nil {
		return a.SessionStorage.GetSession(sessionID)
	}
	return 0, err
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
	a.SessionStorage.DeleteSession(sessionID)

	http.SetCookie(w, &http.Cookie{
		Name:   sessionId,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return nil
}
