package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateUser(ctx context.Context, name, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

func (m *MockService) VerifyUser(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (m *MockService) DeleteUser(ctx context.Context, id int) error {
	return nil
}

func TestHandler_handleCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateUserRequest
		mockSetup      func(*MockService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful user creation",
			requestBody: CreateUserRequest{
				Name:     "test",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(ms *MockService) {
				ms.On("CreateUser", mock.Anything, "test@example.com", "password123").Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "service error",
			requestBody: CreateUserRequest{
				Name:     "test",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(ms *MockService) {
				ms.On("CreateUser", mock.Anything, "test@example.com", "password123").Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "invalid JSON",
			requestBody:    CreateUserRequest{},      // Will be replaced with invalid JSON
			mockSetup:      func(ms *MockService) {}, // No mock calls expected
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			tt.mockSetup(mockService)
			handler := NewHandler(mockService)

			var body []byte
			var err error
			if tt.name == "invalid JSON" {
				body = []byte(`{"invalid": json}`)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.handleCreateUser(w, req, httprouter.Params{})

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)

		})
	}

}
