package authorization

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

func TestAuthorize(t *testing.T) {
	sessionStorage := session.NewSessionStorage()
	sessionUseCase := usecase.NewAuthUseCase(sessionStorage)

	userStorage := user.NewUserStorage()
	userUseCase := usecase.NewUserUseCase(userStorage)

	authHandler := NewAuthorizationHandler(sessionUseCase, userUseCase)

	_, err := userUseCase.CreateUser("san", "A123B123abc")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	authTestCases := []struct {
		name            string
		inputJSON       string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Successful authorization",
			inputJSON:       `{"username": "san", "password": "A123B123abc"}`,
			expectedStatus:  http.StatusOK,
			expectedMessage: "san",
		},
		{
			name:            "Empty username",
			inputJSON:       `{"username": "", "password": "A123B123abc"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Username doesn't meet requirements",
		},
		{
			name:            "Empty password",
			inputJSON:       `{"username": "san", "password": ""}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Password doesn't meet requirements",
		},
		{
			name:            "Simple password",
			inputJSON:       `{"username": "san", "password": "123abc"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Password doesn't meet requirements",
		},
		{
			name:            "User not found",
			inputJSON:       `{"username": "Sanboy", "password": "Abc123abc123"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "User not found",
		},
		{
			name:            "Context without request",
			inputJSON:       `{"username": "san", "password": "A123B123abc"}`,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed getting request",
		},
		{
			name:            "Context without response writer",
			inputJSON:       `{"username": "san", "password": "A123B123abc"}`,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed getting response writer",
		},
	}

	for _, tc := range authTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем запрос
			req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(tc.inputJSON))
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
			if tc.name == "Context request with cookie" {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: "test_session_id",
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
			response, err := authHandler.Authorize(ctx, user)

			// Проверяем ответ
			var httpErr httperrors.HttpError
			if err != nil {
				httpErr = httperrors.UnwrapHttpError(err)
			}

			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedMessage, response.Username, "handler returned wrong username")
			} else {
				assert.Equal(t, tc.expectedMessage, httpErr.Message, "handler returned wrong error message")
				assert.Equal(t, tc.expectedStatus, httpErr.Code, "handler returned wrong status code")
			}
		})
	}

	logoutTestCases := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Successful logout",
			expectedStatus:  0,
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
			expectedStatus:  http.StatusForbidden,
			expectedMessage: "http: named cookie not present",
		},
		{
			name:            "Context with wrong cookie",
			expectedStatus:  http.StatusForbidden,
			expectedMessage: "Session not found",
		},
	}

	for _, tc := range logoutTestCases {
		t.Run(tc.name, func(t *testing.T) {
			sessionStorage.SaveSession("session", 2)
			// Создаем запрос
			req, err := http.NewRequest("POST", "/logout", nil)
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
			response, httpErr := authHandler.LogOut(ctx, entities.User{})

			// Проверяем ответ

			if tc.expectedStatus == 0 {
				assert.Equal(t, tc.expectedMessage, response.Username, "handler returned wrong username")
			} else {
				assert.Equal(t, httpErr.Error(), tc.expectedMessage, "handler returned wrong error message")
			}
		})
	}
}
