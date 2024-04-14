package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"

	sight "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery"
	user "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/cors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

func SetupRouter(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(cors.CorsMiddleware)
	router.Use(middle.SessionMiddleware)

	router.Mount("/sights", SightRoutes())

	// user authorization and registration
	router.Mount("/signup", SignUpRoutes())
	router.Mount("/login", AuthRoutes())
	router.Mount("/logout", LogOutRoutes())

	// user profile
	router.Mount("/profile/{id}", GetProfileRoutes())
	router.Mount("/profile/{id}/edit", EditProfileRoutes())
	router.Mount("/profile/{id}/delete", DeleteProfileRoutes())
	router.Mount("/profile/{id}/reset_password", UpdateUserPasswordRoutes())

	//TODO:нужно приспособить обертку под работу multipart/form-data
	handler := &user.ProfileHandler{}
	router.Post("/profile/{id}/upload", func(w http.ResponseWriter, r *http.Request) {
		handler.UploadFile(w, r)
	})

	// comments
	router.Mount("/sight/{id}", SightByIDRoutes())
	router.Mount("/sight/{id}/create", CreateCommentRoutes())
	router.Mount("/sight/{sid}/edit/{cid}", EditCommentRoutes())
	router.Mount("/sight/{sid}/delete/{cid}", DeleteCommentRoutes())

	//journeys
	router.Mount("/trip/{id}/delete", DeleteJourneyRoutes())
	router.Mount("/trip/create", CreateJourneyRoutes())
	router.Mount("/{userID}/trips", JourneyRoutes())

	// journey_sights
	router.Mount("/trip/{id}", JourneySightRoutes())
	router.Mount("/trip/{id}/sight/add", AddJourneySightRoutes())
	router.Mount("/trip/{id}/sight/delete", DeleteJourneySightRoutes())

	return router
}

func SightRoutes() chi.Router {
	router := chi.NewRouter()
	sightsHandler := sight.SightsHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Sight, entities.Sights]{ServeHTTP: sightsHandler.GetSights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func SignUpRoutes() chi.Router {
	router := chi.NewRouter()

	regHandler := user.RegistrationHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.User, user.UserResponse]{ServeHTTP: regHandler.SignUp}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func LogOutRoutes() chi.Router {
	router := chi.NewRouter()

	logOutHandler := user.AuthorizationHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.User, user.UserResponse]{ServeHTTP: logOutHandler.LogOut}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func AuthRoutes() chi.Router {
	router := chi.NewRouter()

	authHandler := user.AuthorizationHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.User, user.UserResponse]{ServeHTTP: authHandler.Authorize}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func SightByIDRoutes() chi.Router {
	router := chi.NewRouter()
	SightByIDHandler := sight.SightsHandler{}

	wrapperInstance := &wrapper.Wrapper[entities.Sight, sight.SightComments]{ServeHTTP: SightByIDHandler.GetSightByID}

	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func CreateCommentRoutes() chi.Router {
	router := chi.NewRouter()

	commHandler := sight.CommentHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: commHandler.CreateComment}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func EditCommentRoutes() chi.Router {
	router := chi.NewRouter()

	commHandler := sight.CommentHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: commHandler.EditComment}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteCommentRoutes() chi.Router {
	router := chi.NewRouter()

	commHandler := sight.CommentHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Comment, entities.Comment]{ServeHTTP: commHandler.DeleteComment}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func CreateJourneyRoutes() chi.Router {
	router := chi.NewRouter()

	journeyHandler := sight.JourneyHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: journeyHandler.CreateJourney}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteJourneyRoutes() chi.Router {
	router := chi.NewRouter()

	journeyHandler := sight.JourneyHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journey]{ServeHTTP: journeyHandler.DeleteJourney}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func JourneyRoutes() chi.Router {
	router := chi.NewRouter()

	journeyHandler := sight.JourneyHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.Journey, entities.Journeys]{ServeHTTP: journeyHandler.GetJourneys}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func AddJourneySightRoutes() chi.Router {
	router := chi.NewRouter()

	journeyHandler := sight.JourneyHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.JourneySightID, entities.JourneySight]{ServeHTTP: journeyHandler.AddJourneySight}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteJourneySightRoutes() chi.Router {
	router := chi.NewRouter()

	journeyHandler := sight.JourneyHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.JourneySight, entities.JourneySight]{ServeHTTP: journeyHandler.DeleteJourneySight}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func JourneySightRoutes() chi.Router {
	router := chi.NewRouter()

	journeyHandler := sight.JourneyHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.JourneySight, entities.JourneySights]{ServeHTTP: journeyHandler.GetJourneySights}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

// profile
func GetProfileRoutes() chi.Router {
	router := chi.NewRouter()
	profileHandler := user.ProfileHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.User, user.ProfileResponse]{ServeHTTP: profileHandler.GetUserProfile}
	router.Get("/", wrapperInstance.HandlerWrapper)

	return router
}

func DeleteProfileRoutes() chi.Router {
	router := chi.NewRouter()
	profileHandler := user.ProfileHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.User, user.ProfileResponse]{ServeHTTP: profileHandler.DeleteUser}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func EditProfileRoutes() chi.Router {
	router := chi.NewRouter()
	profileHandler := user.ProfileHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.UserProfile, entities.UserProfile]{ServeHTTP: profileHandler.EditUserProfile}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}

func UpdateUserPasswordRoutes() chi.Router {
	router := chi.NewRouter()
	profileHandler := user.ProfileHandler{}
	wrapperInstance := &wrapper.Wrapper[entities.User, user.ProfileResponse]{ServeHTTP: profileHandler.UpdateUserPassword}
	router.Post("/", wrapperInstance.HandlerWrapper)

	return router
}
