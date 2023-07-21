package v1

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/service"
	mock_service "github.com/aidos-dev/habit-tracker/backend/internal/service/mocks"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/handlers/slogdiscard"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func Test_handler_signUpWeb(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user models.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"userName": "testUser",
				"firstName": "testUserName",
				"lastName": "testUserLastName",
				"eMail": "testEmail@gmail.com",
				"password": "qwerty"
				}`,
			inputUser: models.User{
				Username:  "testUser",
				FirstName: "testUserName",
				LastName:  "testUserLastName",
				Email:     "testEmail@gmail.com",
				Password:  "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUser, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputUser)

			log := slogdiscard.NewDiscardLogger()

			services := &service.Service{User: user}
			handler := NewHandler(log, services)

			// Test Server
			r := gin.New()
			r.POST("sign-up", handler.signUpWeb)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST", "/sign-up",
				bytes.NewBufferString(testCase.inputBody),
			)

			// Perform Request
			r.ServeHTTP(w, req)

			if w.Code != testCase.expectedStatusCode {
				t.Errorf("Expected status code: %d but got: %d", testCase.expectedStatusCode, w.Code)
				// t.Error("Status codes don’t match")
				// t.Logf("Expected status code: %d", testCase.expectedStatusCode)
				// t.Logf("But got: %d", w.Code)
			}

			if w.Body.String() != testCase.expectedRequestBody {
				t.Errorf("Expected response body '%s' but got '%s'", testCase.expectedRequestBody, w.Body.String())
			}
		})
	}
}
