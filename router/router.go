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
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/quiz"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/cors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
	httpSwagger "github.com/swaggo/http-swagger"
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
	router.Use(middle.XSSMiddleware)

	// upload image
	router.HandleFunc("/upload", user.Upload)

	router.Mount("/api/sights", SightRoutes(handlers.SightHandler))

	// user authorization and registration
	router.Mount("/api/signup", SignUpRoutes(handlers.RegHandler))
	router.Mount("/api/login", AuthRoutes(handlers.AuthHandler))
	router.Mount("/api/logout", LogOutRoutes(handlers.AuthHandler))

	// user profile
	router.Mount("/api/profile/{id}", GetProfileRoutes(handlers.ProfileHandler))
	router.Mount("/api/profile/{id}/edit", EditProfileRoutes(handlers.ProfileHandler))
	router.Mount("/api/profile/{id}/delete", DeleteProfileRoutes(handlers.DeactivationHandler))
	router.Mount("/api/profile/{id}/reset_password", UpdateUserPasswordRoutes(handlers.AuthHandler))

	//TODO:нужно приспособить обертку под работу multipart/form-data
	router.Post("/profile/{id}/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProfileHandler.UploadFile(w, r)
	})

	// comments
	router.Mount("/api/sight/{id}", GetSightRoutes(handlers.SightHandler))
	router.Mount("/api/sight/{id}/create", CreateCommentRoutes(handlers.CommentHandler))
	router.Mount("/api/sight/{sid}/edit/{cid}", EditCommentRoutes(handlers.CommentHandler))
	router.Mount("/api/sight/{sid}/delete/{cid}", DeleteCommentRoutes(handlers.CommentHandler))
	router.Mount("/api/sight/quiz", SearchSightsRoutes(handlers.SightHandler))

	//journeys
	router.Mount("/api/trip/{id}/delete", DeleteJourneyRoutes(handlers.JourneyHandler))
	router.Mount("/api/trip/create", CreateJourneyRoutes(handlers.JourneyHandler))
	router.Mount("/api/{userID}/trips", JourneyRoutes(handlers.JourneyHandler))

	// journey_sights
	router.Mount("/api/trip/{id}", JourneySightRoutes(handlers.JourneyHandler))
	router.Mount("/api/trip/{id}/sight/add", AddJourneySightRoutes(handlers.JourneyHandler))
	router.Mount("/api/trip/{id}/edit", EditJourney(handlers.JourneyHandler))
	router.Mount("/api/trip/{id}/sight/delete", DeleteJourneySightRoutes(handlers.JourneyHandler))

	// quiz
	router.Mount("/api/review/create", CreateReviewRoutes(handlers.QuizHandler))
	router.Mount("/api/review/check", CheckUserReviewRoutes(handlers.QuizHandler))
	router.Mount("/api/review/get", GetStatistic(handlers.QuizHandler))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

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

func EditJourney(handler *journey.JourneyHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: handler.EditJourney}
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

func SearchSightsRoutes(handler *sight.SightHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Sight, entities.Sights]{ServeHTTP: handler.SearchSights}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func CreateReviewRoutes(handler *quiz.QuizHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Review, bool]{ServeHTTP: handler.CreateReview}
	router.Post("/", wrapperInstance.HandlerWrapper)
	return router
}

func CheckUserReviewRoutes(handler *quiz.QuizHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Review, entities.DataCheck]{ServeHTTP: handler.CheckData}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}

func GetStatistic(handler *quiz.QuizHandler) chi.Router {
	router := chi.NewRouter()
	wrapperInstance := &wrapper.Wrapper[entities.Statistic, []entities.Statistic]{ServeHTTP: handler.SetStat}
	router.Get("/", wrapperInstance.HandlerWrapper)
	return router
}
