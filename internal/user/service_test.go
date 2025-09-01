package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
func (m *MockRepo) DeleteUser(ctx context.Context, id int) error {
	return nil
}

func TestService_handleCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		user           *User
		mockSetup      func(*MockRepo)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful user creation",
			user: &User{
				Name:     "test",
				Email:    "test@test.com",
				Password: "password",
			},
			mockSetup: func(ms *MockRepo) {
				ms.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tt.mockSetup(mockRepo)
			service := NewService(mockRepo)

			err := service.CreateUser(context.Background(), tt.user.Name, tt.user.Email, tt.user.Password)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
