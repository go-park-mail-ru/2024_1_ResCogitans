package registration

import (
	"bytes"
	"encoding/json"
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

func TestSignUp(t *testing.T) {
	sessionStorage := session.NewSessionStorage()
	sessionUseCase := usecase.NewAuthUseCase(sessionStorage)

	userStorage := user.NewUserStorage()
	userUseCase := usecase.NewUserUseCase(userStorage)

	signUpHandler := NewRegistrationHandler(sessionUseCase, userUseCase)

	testCases := []struct {
		name            string
		inputJSON       string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Successful registration",
			inputJSON:       `{"username": "san", "password": "A123B123abc"}`,
			expectedStatus:  http.StatusOK,
			expectedMessage: "san",
		},
		{
			name:            "Username taken",
			inputJSON:       `{"username": "san", "password": "A123B123abc123"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Username is taken",
		},
		{
			name:            "Simple username",
			inputJSON:       `{"username": "US", "password": "A123B123abc123"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Username doesn't meet requirements",
		},
		{
			name:            "Simple password",
			inputJSON:       `{"username": "sanBoy", "password": "1234"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Password doesn't meet requirements",
		},
		{
			name:            "Context without response writer",
			inputJSON:       `{"username": "sanBoy", "password": "A123B123abc123"}`,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed getting response writer",
		},
		{
			name:            "Context without request",
			inputJSON:       `{"username": "sanBoy", "password": "A123B123abc123"}`,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed getting request",
		},
		{
			name:            "Request with cookie",
			inputJSON:       `{"username": "sancho", "password": "A123B123abc123"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "User is already authorized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем запрос
			sessionStorage.SaveSession("session", 2)
			req, err := http.NewRequest("POST", "/signup", bytes.NewBufferString(tc.inputJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Создаем контекст с дополнительными данными
			ctx := req.Context()
			if tc.name != "Context without request" {
				ctx = httputils.SetRequestToCtx(ctx, req)
			}
			if tc.name != "Context without response writer" {
				ctx = httputils.SetResponseWriterToCtx(ctx, httptest.NewRecorder())
			}
			if tc.name == "Request with cookie" {
				encoded, err := usecase.CookieHandler.Encode("session_id", "session")
				if err != nil {
					t.Fatal(err)
				}
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: encoded,
				})
			} else {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: "",
				})
			}

			// Создаем пользователя и декодируем JSON
			user := entities.User{}
			err = json.NewDecoder(req.Body).Decode(&user)
			if err != nil {
				t.Fatal(err)
			}

			// Вызываем ручку с контекстом
			response, err := signUpHandler.SignUp(ctx, user)
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
