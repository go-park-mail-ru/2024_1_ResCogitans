package delivery_test

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"strings"
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
// 	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
// 	"github.com/stretchr/testify/assert"
// )

// func TestSignUp(t *testing.T) {
// 	regHandler := registration.Registration{}

// 	testCases := []struct {
// 		name           string
// 		inputJSON      string
// 		expectedStatus int
// 		expectedError  string
// 	}{
// 		{
// 			name:           "Successful registration",
// 			inputJSON:      `{"username": "testuser", "password": "testpassword"}`,
// 			expectedStatus: http.StatusCreated,
// 			expectedError:  "",
// 		},
// 		{
// 			name:           "Empty username",
// 			inputJSON:      `{"username": "", "password": "testpassword"}`,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  "username and password must not be empty",
// 		},
// 		{
// 			name:           "Empty password",
// 			inputJSON:      `{"username": "testuser", "password": ""}`,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  "username and password must not be empty",
// 		},
// 		{
// 			name:           "User already exists",
// 			inputJSON:      `{"username": "testuser", "password": "testpassword"}`,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedError:  "username already exists",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			req, err := http.NewRequest("POST", "/signup", strings.NewReader(tc.inputJSON))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			var user entities.User
// 			err = json.NewDecoder(req.Body).Decode(&user)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			ctx := context.WithValue(req.Context(), "requestData", user)

// 			response, err := regHandler.SignUp(ctx)
// 			if err != nil && tc.expectedError == "" {
// 				t.Errorf("Ошибка при вызове SignUp: %v", err)
// 			}

// 			assert.Equal(t, tc.expectedStatus, response.Status, "handler returned wrong status code")

// 			if tc.expectedError != "" {
// 				assert.EqualError(t, err, tc.expectedError, "handler returned unexpected error message")
// 			}
// 		})
// 	}
// }
