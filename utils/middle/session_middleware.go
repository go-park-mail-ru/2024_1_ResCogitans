package middle

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
)

type AuthMiddleware struct {
	s *usecase.SessionUseCase
	c *usecase.CSRFUseCase
}

func NewAuthMiddleware(session *usecase.SessionUseCase, csrf *usecase.CSRFUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		s: session,
		c: csrf,
	}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.s.GetSession(r.Context(), r)
		if !httperrors.IsHttpError(err) && err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if userID != 0 && r.Method == http.MethodPost {
			token := r.Header.Get("X-CSRF-Token")
			if token == "" {
				newToken, err := m.c.CreateToken(r.Context(), userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				m.c.SetToken(newToken, w)
			} else {
				err = m.c.CompareToken(r.Context(), token, userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				}
				newToken, err := m.c.UpdateToken(r.Context(), userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				m.c.SetToken(newToken, w)
			}
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
