package middle

import (
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := bluemonday.UGCPolicy()

		// Фильтрация каждого элемента формы
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for k, v := range r.Form {
			filtered := make([]string, len(v))
			for i, value := range v {
				filtered[i] = p.Sanitize(value)
			}
			r.Form[k] = filtered
		}

		// Фильтрация каждого элемента заголовка
		for k, v := range r.Header {
			filtered := make([]string, len(v))
			for i, value := range v {
				filtered[i] = p.Sanitize(value)
			}
			r.Header[k] = filtered
		}

		next.ServeHTTP(w, r)
	})
}
