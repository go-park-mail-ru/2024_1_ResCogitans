package updateUserData

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

func TestUpdate(t *testing.T) {
	sessionStorage := session.NewSessionStorage()
	sessionUseCase := usecase.NewAuthUseCase(sessionStorage)

	userStorage := user.NewUserStorage()
	userUseCase := usecase.NewUserUseCase(userStorage)

	updateHandler := NewUpdateDataHandler(sessionUseCase, userUseCase)

	user, err := userUseCase.CreateUser("san", "A123B123abc")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	testCases := []struct {
		name            string
		inputJSON       string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Successful data update",
			inputJSON:       `{"username": "sanBoy", "password": "A123B1234abc"}`,
			expectedStatus:  http.StatusOK,
			expectedMessage: "sanBoy",
		},
		{
			name:            "Context without request",
			inputJSON:       `{"username": "san", "password": "A123B123"}`,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed getting request",
		},
		{
			name:            "Request without session",
			inputJSON:       `{"username": "san", "password": "A123B123"}`,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Error decoding cookie",
		},
		{
			name:            "Request with wrong session",
			inputJSON:       `{"username": "san", "password": "A123B123"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Session not found",
		},
		{
			name:            "Simple password",
			inputJSON:       `{"username": "san", "password": "A123"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Password doesn't meet requirements",
		},
		{
			name:            "Simple username",
			inputJSON:       `{"username": "US", "password": "A123B123"}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Username doesn't meet requirements",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем запрос
			sessionStorage.SaveSession("session", user.ID)
			req, err := http.NewRequest("POST", "/updatedata", bytes.NewBufferString(tc.inputJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Создаем контекст с дополнительными данными
			ctx := req.Context()
			if tc.name != "Context without request" {
				ctx = httputils.SetRequestToCtx(ctx, req)
			}
			ctx = httputils.SetResponseWriterToCtx(ctx, httptest.NewRecorder())
			if tc.name != "Request without session" {
				var encoded string
				if tc.name == "Request with wrong session" {
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
			response, err := updateHandler.Update(ctx, user)
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
