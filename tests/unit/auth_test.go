package unit_test

import (
	"context"
	"testing"
	"time"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/token"
	"shifty-backend/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Helper tạo Token dùng chung
func newTestTokenMaster() *token.TokenMaster {
	return token.NewToken(
		"secret_access_test_key",
		"secret_refresh_test_key",
		time.Minute*15,
		time.Hour*24,
	)
}
func TestAuthUseCase_RegisterLocal(t *testing.T) {
	mockUser := &entity.User{
		Email:    "newuser@gmail.com",
		Password: "hashedpassword",
		Role:     "User",
	}

	tests := []struct {
		name          string
		user          *entity.User
		mockSetup     func(*MockUserRepo)
		expectedError bool
	}{
		{
			name: "Success - Register New User",
			user: mockUser,
			mockSetup: func(u *MockUserRepo) {
				u.On("IsEmailExist", mock.Anything, mockUser.Email).Return(false, nil)
				u.On("Create", mock.Anything, mockUser).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Fail - Email Already Exists",
			user: mockUser,
			mockSetup: func(u *MockUserRepo) {
				u.On("IsEmailExist", mock.Anything, mockUser.Email).Return(true, nil)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepo)
			mockRedisRepo := new(MockRedisRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockUserRepo)
			}

			authUC := usecase.NewAuthUseCase(mockUserRepo, newTestTokenMaster(), time.Second, mockRedisRepo, nil, nil)
			err := authUC.RegisterLocal(context.Background(), tt.user)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_LoginLocal(t *testing.T) {
	pass := "test123@"
	hashedPass, _ := utils.HashPassword(pass)
	userID := uuid.New()

	mockUser := &entity.User{
		ID:          userID,
		Email:       "test@gmail.com",
		Password:    hashedPass,
		AccountType: "Local",
		Role:        "User",
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(*MockUserRepo, *MockRedisRepo)
		expectedError bool
	}{
		{
			name:     "Success - Login OK",
			email:    "test@gmail.com",
			password: pass,
			mockSetup: func(u *MockUserRepo, r *MockRedisRepo) {
				u.On("GetByEmail", mock.Anything, "test@gmail.com").Return(mockUser, nil)
				r.On("CreateSession", mock.Anything, mock.Anything).Return(nil)
				r.On("SaveUserCache", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "Fail - User Not Found",
			email:    "ghost@gmail.com",
			password: pass,
			mockSetup: func(u *MockUserRepo, r *MockRedisRepo) {
				u.On("GetByEmail", mock.Anything, "ghost@gmail.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepo)
			mockRedisRepo := new(MockRedisRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockUserRepo, mockRedisRepo)
			}

			authUC := usecase.NewAuthUseCase(mockUserRepo, newTestTokenMaster(), time.Second, mockRedisRepo, nil, nil)
			_, _, _, err := authUC.LoginLocal(context.Background(), tt.email, tt.password, "UA", "IP")

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockUserRepo.AssertExpectations(t)
			mockRedisRepo.AssertExpectations(t)
		})
	}
}
