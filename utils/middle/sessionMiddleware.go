package middle

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"net/http"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := entities.GetSession(r)

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
