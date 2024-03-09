package middle

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
)

type userIDType struct{}

var userIDKey userIDType

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := usecase.GetSession(r)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
