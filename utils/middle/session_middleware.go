package middle

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
)

type AuthMiddleware struct {
	useCase usecase.SessionInterface
}

func NewAuthMiddleware(useCase usecase.SessionInterface) *AuthMiddleware {
	return &AuthMiddleware{
		useCase: useCase,
	}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := m.useCase.GetSession(r)

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
