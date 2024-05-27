package usecase

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/service/gen"
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
	client gen.SessionServiceClient
}

func NewSessionUseCase(conn *grpc.ClientConn) *SessionUseCase {
	return &SessionUseCase{
		client: gen.NewSessionServiceClient(conn),
	}
}

func (s *SessionUseCase) CreateSession(ctx context.Context, w http.ResponseWriter, userID int) error {
	sessionID := uuid.New().String()
	response, err := s.client.CreateSession(ctx, &gen.SaveSessionRequest{
		SessionID: sessionID,
		UserID:    int32(userID),
	})
	if response.Error != "" {
		return errors.New(response.Error)
	}
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

func (s *SessionUseCase) GetSession(ctx context.Context, r *http.Request) (int, error) {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return 0, httperrors.NewHttpError(http.StatusUnauthorized, "Cookie not found")
		}
		return 0, err
	}

	var sessionID string
	if err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID); err == nil {
		ans, err := s.client.GetSession(ctx, &gen.GetSessionRequest{SessionID: sessionID})
		if err != nil {
			return 0, err
		}
		return int(ans.UserID), nil
	}
	return 0, httperrors.NewHttpError(http.StatusInternalServerError, "Error decoding cookie")
}

func (s *SessionUseCase) ClearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return err
	}

	var sessionID string
	err = CookieHandler.Decode(sessionId, cookie.Value, &sessionID)
	if err != nil {
		return err
	}
	_, err = s.client.DeleteSession(ctx, &gen.DeleteSessionRequest{SessionID: sessionID})
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
