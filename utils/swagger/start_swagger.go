package swagger

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func StartSwagger() error {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"),
	))

	go func() {
		err := http.ListenAndServe(":1323", r)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}
