package unit_test

import (
	"context"
	"errors"
	"testing"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)



func TestUserUseCase_FindUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockSetup     func(*MockUserRepo)
		expectedError bool
	}{
		{
			name:  "Success - Found User",
			email: "test@gmail.com",
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByEmail", mock.Anything, "test@gmail.com").
					Return(&entity.User{Email: "test@gmail.com"}, nil)
			},
			expectedError: false,
		},
		{
			name:  "Fail - User Not Found",
			email: "ghost@gmail.com",
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByEmail", mock.Anything, "ghost@gmail.com").
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
		},
		{
			name:  "Fail - Database Error",
			email: "db@gmail.com",
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByEmail", mock.Anything, "db@gmail.com").
					Return(nil, errors.New("db error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			mockUploader := new(MockUploader)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			userUC := usecase.NewUserUseCase(mockRepo, mockUserResRepo, mockUploader)
			user, err := userUC.FindUserByEmail(context.Background(), tt.email)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_FindUserByID(t *testing.T) {
	targetID := uuid.New().String()

	tests := []struct {
		name          string
		id            string
		mockSetup     func(*MockUserRepo)
		expectedError bool
	}{
		{
			name: "Success - Found ID",
			id:   targetID,
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByID", mock.Anything, targetID).
					Return(&entity.User{FullName: "Found"}, nil)
			},
			expectedError: false,
		},
		{
			name: "Fail - Not Found",
			id:   "bad-id",
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByID", mock.Anything, "bad-id").
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			mockUploader := new(MockUploader)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			userUC := usecase.NewUserUseCase(mockRepo, mockUserResRepo, mockUploader)
			_, err := userUC.FindUserByID(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
