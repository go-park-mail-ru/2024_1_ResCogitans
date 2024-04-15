package usecase

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/session"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"
)

const sessionId = "session_id"

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type SessionInterface interface {
	CreateSession(w http.ResponseWriter, userID int) error
	GetSession(r *http.Request) (int, error)
	ClearSession(w http.ResponseWriter, r *http.Request) error
}

type SessionUseCase struct {
	SessionStorage session.StorageInterface
}

func NewSessionUseCase(storage session.StorageInterface) SessionInterface {
	return &SessionUseCase{
		SessionStorage: storage,
	}
}

func (a *SessionUseCase) CreateSession(w http.ResponseWriter, userID int) error {
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

func (a *SessionUseCase) GetSession(r *http.Request) (int, error) {
	println("11")
	cookie, err := r.Cookie(sessionId)
	if errors.Is(err, http.ErrNoCookie) {
		return 0, httperrors.NewHttpError(http.StatusInternalServerError, "Cookie not found")
	}

	if err != nil {
		return 0, err
	}

	println("22")

	var sessionID string
	if err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID); err == nil {
		return a.SessionStorage.GetSession(sessionID)
	}
	return 0, httperrors.NewHttpError(http.StatusInternalServerError, "Error decoding cookie")
}

func (a *SessionUseCase) ClearSession(w http.ResponseWriter, r *http.Request) error {
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
