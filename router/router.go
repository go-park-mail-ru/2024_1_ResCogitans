package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/authorization"
	user "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/avatar"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/comment"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/deactivation"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/journey"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/profile"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/quiz"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/cors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
	"golang.org/x/exp/slog"
)

func SetupRouter(logger *slog.Logger, handlers *initialization.Handlers) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(cors.CorsMiddleware)
	router.Use(handlers.AuthMiddleware.Auth)
	router.Use(middle.XSSMiddleware)

	// upload image
	router.HandleFunc("/upload", user.Upload)

	router.Mount("/api/sights", SightRoutes(handlers.SightHandler, logger))

	// user authorization and registration
	router.Mount("/api/signup", SignUpRoutes(handlers.RegHandler, logger))
	router.Mount("/api/login", AuthRoutes(handlers.AuthHandler, logger))
	router.Mount("/api/logout", LogOutRoutes(handlers.AuthHandler, logger))

	// user profile
	router.Mount("/api/profile/{id}", GetProfileRoutes(handlers.ProfileHandler, logger))
	router.Mount("/api/profile/{id}/edit", EditProfileRoutes(handlers.ProfileHandler, logger))
	router.Mount("/api/profile/{id}/delete", DeleteProfileRoutes(handlers.DeactivationHandler, logger))
	router.Mount("/api/profile/{id}/reset_password", UpdateUserPasswordRoutes(handlers.AuthHandler, logger))

	//TODO:нужно приспособить обертку под работу multipart/form-data
	router.Post("/profile/{id}/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProfileHandler.UploadFile(w, r, logger)
	})

	// comments
	router.Mount("/api/sight/{id}", GetSightRoutes(handlers.SightHandler, logger))
	router.Mount("/api/sight/{id}/create", CreateCommentRoutes(handlers.CommentHandler, logger))
	router.Mount("/api/sight/{sid}/edit/{cid}", EditCommentRoutes(handlers.CommentHandler, logger))
	router.Mount("/api/sight/{sid}/delete/{cid}", DeleteCommentRoutes(handlers.CommentHandler, logger))
	router.Mount("/api/sight/quiz", SearchSightsRoutes(handlers.SightHandler, logger))

	//journeys
	router.Mount("/api/trip/{id}/delete", DeleteJourneyRoutes(handlers.JourneyHandler, logger))
	router.Mount("/api/trip/create", CreateJourneyRoutes(handlers.JourneyHandler, logger))
	router.Mount("/api/{userID}/trips", JourneyRoutes(handlers.JourneyHandler, logger))

	// journey_sights
	router.Mount("/api/trip/{id}", JourneySightRoutes(handlers.JourneyHandler, logger))
	router.Mount("/api/trip/{id}/sight/add", AddJourneySightRoutes(handlers.JourneyHandler, logger))
	router.Mount("/api/trip/{id}/edit", EditJourney(handlers.JourneyHandler, logger))
	router.Mount("/api/trip/{id}/sight/delete", DeleteJourneySightRoutes(handlers.JourneyHandler, logger))

	// quiz
	router.Mount("/api/review/create", CreateReviewRoutes(handlers.QuizHandler, logger))
	router.Mount("/api/review/check", CheckUserReviewRoutes(handlers.QuizHandler, logger))
	router.Mount("/api/review/get", GetStatistic(handlers.QuizHandler, logger))

	return router
}

func SightRoutes(handler *sight.SightHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Sight, entities.Sights]{ServeHTTP: handler.GetSights, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func SignUpRoutes(handler *registration.RegistrationHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.SignUp, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func LogOutRoutes(handler *authorization.AuthorizationHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.LogOut, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func AuthRoutes(handler *authorization.AuthorizationHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.Authorize, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func GetSightRoutes(handler *sight.SightHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Sight, entities.SightComments]{ServeHTTP: handler.GetSight, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func CreateCommentRoutes(handler *comment.CommentHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: handler.CreateComment, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func EditCommentRoutes(handler *comment.CommentHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: handler.EditComment, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func DeleteCommentRoutes(handler *comment.CommentHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: handler.DeleteComment, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func CreateJourneyRoutes(handler *journey.JourneyHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: handler.CreateJourney, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func DeleteJourneyRoutes(handler *journey.JourneyHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: handler.DeleteJourney, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func JourneyRoutes(handler *journey.JourneyHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journeys]{ServeHTTP: handler.GetJourneys, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func AddJourneySightRoutes(handler *journey.JourneyHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.JourneySightID, entities.JourneySight]{ServeHTTP: handler.AddJourneySight, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func EditJourney(handler *journey.JourneyHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: handler.EditJourney, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func DeleteJourneySightRoutes(handler *journey.JourneyHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.JourneySight, entities.JourneySight]{ServeHTTP: handler.DeleteJourneySight, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func JourneySightRoutes(handler *journey.JourneyHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.JourneySight, entities.JourneySights]{ServeHTTP: handler.GetJourneySights, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

// profile
func GetProfileRoutes(handler *profile.ProfileHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.UserProfile, entities.UserProfile]{ServeHTTP: handler.Get, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func EditProfileRoutes(handler *profile.ProfileHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.UserProfile, entities.UserProfile]{ServeHTTP: handler.Edit, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func DeleteProfileRoutes(handler *deactivation.DeactivationHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.Deactivate, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func UpdateUserPasswordRoutes(handler *authorization.AuthorizationHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.User, entities.UserResponse]{ServeHTTP: handler.UpdatePassword, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func SearchSightsRoutes(handler *sight.SightHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Sight, entities.Sights]{ServeHTTP: handler.SearchSights, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func CreateReviewRoutes(handler *quiz.QuizHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Review, bool]{ServeHTTP: handler.CreateReview, Logger: logger}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func CheckUserReviewRoutes(handler *quiz.QuizHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Review, entities.DataCheck]{ServeHTTP: handler.CheckData, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func GetStatistic(handler *quiz.QuizHandler, logger *slog.Logger) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Statistic, []entities.Statistic]{ServeHTTP: handler.SetStat, Logger: logger}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}
