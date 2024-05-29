package usecase

import (
	"context"
	"net/http"
	"time"

	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
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
	CreateSession(ctx context.Context, w http.ResponseWriter, userID int) error
	GetSession(ctx context.Context, r *http.Request) (int, error)
	ClearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type SessionUseCase struct {
	SessionStorage storage.SessionStorageInterface
}

func NewSessionUseCase(storage storage.SessionStorageInterface) *SessionUseCase {
	return &SessionUseCase{
		SessionStorage: storage,
	}
}

func (a *SessionUseCase) CreateSession(ctx context.Context, w http.ResponseWriter, userID int) error {
	sessionID := uuid.New().String()
	err := a.SessionStorage.SaveSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}
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

func (a *SessionUseCase) GetSession(ctx context.Context, r *http.Request) (int, error) {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return 0, httperrors.NewHttpError(http.StatusUnauthorized, "cookie not found")
		}
		return 0, err
	}

	var sessionID string
	if err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID); err == nil {
		return a.SessionStorage.GetSession(ctx, sessionID)
	}
	return 0, errors.New("error decoding cookie")
}

func (a *SessionUseCase) ClearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return err
	}

	var sessionID string
	err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID)
	if err != nil {
		return err
	}
	err = a.SessionStorage.DeleteSession(ctx, sessionID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:   sessionId,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return nil
}
