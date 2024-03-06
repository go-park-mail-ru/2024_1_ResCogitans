package login_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/login"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	authHandler := login.Authorization{}

	entities.CreateUser("san", "123")

	testCases := []struct {
		name           string
		inputJSON      string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Successful authorization",
			inputJSON:      `{"username": "san", "password": "123"}`,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Empty username",
			inputJSON:      `{"username": "", "password": "testpassword"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization failed",
		},
		{
			name:           "Empty password",
			inputJSON:      `{"username": "testuser", "password": ""}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization failed",
		},
		{
			name:           "User not found",
			inputJSON:      `{"username": "nonexistentuser", "password": "testpassword"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/login", strings.NewReader(tc.inputJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			var user entities.User
			err = json.NewDecoder(req.Body).Decode(&user)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), "responseWriter", rr)
			ctx = context.WithValue(ctx, "requestData", user)

			response, err := authHandler.Authorize(ctx)

			assert.Equal(t, tc.expectedStatus, response.Status, "handler returned wrong status code")
		})
	}
}
