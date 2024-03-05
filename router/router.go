package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/login"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

func SetupRouter(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	router.Mount("/sights", SightRoutes())
	router.Mount("/signup", SignUpRoutes())
	router.Mount("/login", AuthRoutes())

	return router
}

func SightRoutes() chi.Router {
	router := chi.NewRouter()
	sightsHandler := sight.SightsHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Sight, sight.Sights]{ServeHTTP: sightsHandler.GetSights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router

}

func SignUpRoutes() chi.Router {
	router := chi.NewRouter()

	regHandler := registration.Registration{}
	wrapperInstance := &wrapper.Wrapper[entities.User, *entities.User]{ServeHTTP: regHandler.SignUp}
	router.Post("/", wrapperInstance.HandlerWrapper)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Signup Page</h1>")
		fmt.Fprintln(w, "<form method='post' action='/signup'>")
		fmt.Fprintln(w, "<label>Username:</label>")
		fmt.Fprintln(w, "<input type='text' name='username'>")
		fmt.Fprintln(w, "<label>Password:</label>")
		fmt.Fprintln(w, "<input type='password' name='password'>")
		fmt.Fprintln(w, "<button type='submit'>Sign Up</button>")
		fmt.Fprintln(w, "</form>")
	})

	return router
}

func AuthRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Login Page</h1>")
		fmt.Fprintln(w, "<form method='post' action='/login'>")
		fmt.Fprintln(w, "<label>Username:</label>")
		fmt.Fprintln(w, "<input type='text' name='username'>")
		fmt.Fprintln(w, "<label>Password:</label>")
		fmt.Fprintln(w, "<input type='password' name='password'>")
		fmt.Fprintln(w, "<button type='submit'>Login</button>")
		fmt.Fprintln(w, "</form>")
	})

	authHandler := login.Authorization{}
	wrapperInstance := &wrapper.Wrapper[entities.User, *entities.User]{ServeHTTP: authHandler.Authorize}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}
