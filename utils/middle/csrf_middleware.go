package middle

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/pkg/errors"
)

type CookieMiddleware struct {
	useCase usecase.CookieInterface
}

func NewCookieMiddleware(useCase usecase.CookieInterface) *CookieMiddleware {
	return &CookieMiddleware{
		useCase: useCase,
	}
}

func (m *CookieMiddleware) CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := m.useCase.GetCookie(r); err != nil {
			if !errors.Is(err, http.ErrNoCookie) || err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				err = m.useCase.Set(w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		if r.Method == http.MethodPost {
			if err := m.useCase.CompareCSRF(r); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			csrfToken, err := m.useCase.ChangeCSRF(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("X-CSRF-Token", csrfToken)
		}

		next.ServeHTTP(w, r)
	})
}
