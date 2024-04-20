package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/authorization"
	user "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/avatar"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/comment"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/deactivation"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/journey"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/profile"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/cors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

func SetupRouter(_ *config.Config, handlers *initialization.Handlers) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(cors.CorsMiddleware)
	router.Use(handlers.AuthMiddleware.Auth)

	// upload image
	router.HandleFunc("/upload", user.Upload)

	router.Mount("/sights", SightRoutes(handlers.SightHandler))

	// user authorization and registration
	router.Mount("/signup", SignUpRoutes(handlers.RegHandler))
	router.Mount("/login", AuthRoutes(handlers.AuthHandler))
	router.Mount("/logout", LogOutRoutes(handlers.AuthHandler))

	// user profile
	router.Mount("/profile/{id}", GetProfileRoutes(handlers.ProfileHandler))
	router.Mount("/profile/{id}/edit", EditProfileRoutes(handlers.ProfileHandler))
	router.Mount("/profile/{id}/delete", DeleteProfileRoutes(handlers.DeactivationHandler))
	router.Mount("/profile/{id}/reset_password", UpdateUserPasswordRoutes(handlers.AuthHandler))

	//TODO:нужно приспособить обертку под работу multipart/form-data
	router.Post("/profile/{id}/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProfileHandler.UploadFile(w, r)
	})

	// comments
	router.Mount("/sight/{id}", GetSightRoutes(handlers.SightHandler))
	router.Mount("/sight/{id}/create", CreateCommentRoutes(handlers.CommentHandler))
	router.Mount("/sight/{sid}/edit/{cid}", EditCommentRoutes(handlers.CommentHandler))
	router.Mount("/sight/{sid}/delete/{cid}", DeleteCommentRoutes(handlers.CommentHandler))

	//journeys
	router.Mount("/trip/{id}/delete", DeleteJourneyRoutes(handlers.JourneyHandler))
	router.Mount("/trip/create", CreateJourneyRoutes(handlers.JourneyHandler))
	router.Mount("/{userID}/trips", JourneyRoutes(handlers.JourneyHandler))

	// journey_sights
	router.Mount("/trip/{id}", JourneySightRoutes(handlers.JourneyHandler))
	router.Mount("/trip/{id}/sight/add", AddJourneySightRoutes(handlers.JourneyHandler))
	router.Mount("/trip/{id}/sight/delete", DeleteJourneySightRoutes(handlers.JourneyHandler))

	return router
}

func SightRoutes(handler *sight.SightHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Sight, entities.Sights]{ServeHTTP: handler.GetSights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func SignUpRoutes(handler *registration.RegistrationHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.SignUp}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func LogOutRoutes(handler *authorization.AuthorizationHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.LogOut}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func AuthRoutes(handler *authorization.AuthorizationHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.Authorize}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func GetSightRoutes(handler *sight.SightHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.Sight, entities.SightComments]{ServeHTTP: handler.GetSight}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func CreateCommentRoutes(handler *comment.CommentHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: handler.CreateComment}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func EditCommentRoutes(handler *comment.CommentHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: handler.EditComment}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteCommentRoutes(handler *comment.CommentHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: handler.DeleteComment}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func CreateJourneyRoutes(handler *journey.JourneyHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: handler.CreateJourney}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteJourneyRoutes(handler *journey.JourneyHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: handler.DeleteJourney}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func JourneyRoutes(handler *journey.JourneyHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journeys]{ServeHTTP: handler.GetJourneys}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func AddJourneySightRoutes(handler *journey.JourneyHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.JourneySightID, entities.JourneySight]{ServeHTTP: handler.AddJourneySight}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteJourneySightRoutes(handler *journey.JourneyHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.JourneySight, entities.JourneySight]{ServeHTTP: handler.DeleteJourneySight}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func JourneySightRoutes(handler *journey.JourneyHandler) chi.Router {
	router := chi.NewRouter()

	wrapperInstance := &wrapper.Wrapper[entities.JourneySight, entities.JourneySights]{ServeHTTP: handler.GetJourneySights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

// profile
func GetProfileRoutes(handler *profile.ProfileHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.UserProfile, entities.UserProfile]{ServeHTTP: handler.Get}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func EditProfileRoutes(handler *profile.ProfileHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.UserProfile, entities.UserProfile]{ServeHTTP: handler.Edit}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteProfileRoutes(handler *deactivation.DeactivationHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.Deactivate}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func UpdateUserPasswordRoutes(handler *authorization.AuthorizationHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.UpdatePassword}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}
