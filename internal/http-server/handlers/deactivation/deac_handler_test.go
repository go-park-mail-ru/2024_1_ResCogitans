package deactivation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/session"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/user"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/stretchr/testify/assert"
)

func TestDeactivate(t *testing.T) {
	sessionStorage := session.NewSessionStorage()
	sessionUseCase := usecase.NewAuthUseCase(sessionStorage)

	userStorage := user.NewUserStorage()
	userUseCase := usecase.NewUserUseCase(userStorage)

	deactivateHandler := NewDeactivationHandler(sessionUseCase, userUseCase)

	testCases := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Successful deactivation",
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
		{
			name:            "Context without request",
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed getting request",
		},
		{
			name:            "Context without response writer",
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed getting response writer",
		},
		{
			name:            "Context without cookie",
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "http: named cookie not present",
		},
		{
			name:            "Context with wrong cookie",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Session not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := userUseCase.CreateUser("san", "A123B123abc")
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}
			sessionStorage.SaveSession("session", user.ID)
			// Создаем запрос
			req, err := http.NewRequest("POST", "/deactivate", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Создаем контекст с дополнительными данными
			ctx := req.Context()
			if tc.name != "Context without request" {
				ctx = httputils.SetRequestToCtx(ctx, req)
			}
			if tc.name != "Context without response writer" {
				ctx = httputils.SetResponseWriterToCtx(ctx, httptest.NewRecorder())
			}
			if tc.name != "Context without cookie" {
				var encoded string
				if tc.name == "Context with wrong cookie" {
					encoded, err = usecase.CookieHandler.Encode("session_id", "session1")
					if err != nil {
						t.Fatal(err)
					}
				} else {
					encoded, err = usecase.CookieHandler.Encode("session_id", "session")
					if err != nil {
						t.Fatal(err)
					}
				}
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: encoded,
				})
			}

			// Вызываем ручку с контекстом
			response, err := deactivateHandler.Deactivate(ctx, entities.User{})

			var httpErr httperrors.HttpError
			if err != nil {
				httpErr = httperrors.UnwrapHttpError(err)
			}

			// Проверяем ответ

			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedMessage, response.Username, "handler returned wrong username")
			} else {
				assert.Equal(t, tc.expectedStatus, httpErr.Code, "handler returned wrong status code")
				assert.Equal(t, tc.expectedMessage, httpErr.Message, "handler returned wrong error message")
			}
		})
	}
}
