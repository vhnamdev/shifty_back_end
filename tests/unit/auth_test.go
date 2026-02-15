package unit_test

import (
	"context"
	"testing"
	"time"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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
				u.On("Create", mock.Anything, mockUser).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Fail - Database Error (Duplicate Email)",
			user: mockUser,
			mockSetup: func(u *MockUserRepo) {
				u.On("Create", mock.Anything, mockUser).Return(xerror.Internal("Database error"))
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

			// Khởi tạo UseCase
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
	posID := uuid.New()
	resID := uuid.New()
	mockUser := &entity.User{
		ID:           userID,
		Email:        "test@gmail.com",
		Password:     hashedPass,
		AccountType:  "Local",
		Role:         "User",
		GoogleID:     "123",
		PositionID:   &posID,
		RestaurantID: &resID,
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
		{
			name:     "Fail - Wrong Password",
			email:    "test@gmail.com",
			password: "wrong_password",
			mockSetup: func(u *MockUserRepo, r *MockRedisRepo) {
				u.On("GetByEmail", mock.Anything, "test@gmail.com").Return(mockUser, nil)
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
			tm := newTestTokenMaster()
			authUC := usecase.NewAuthUseCase(
				mockUserRepo,
				tm,
				time.Second,
				mockRedisRepo,
				nil,
				nil,
			)

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

func TestAuthUseCase_SaveOTP(t *testing.T) {
	email := "test@gmail.com"
	otp := "12345"
	purpose := constants.PurposeRegister

	tests := []struct {
		name          string
		email         string
		otp           string
		purpose       constants.OTPPurpose
		mockSetup     func(*MockRedisRepo)
		expectedError bool
	}{
		{
			name:    "Success - Save OTP",
			email:   email,
			otp:     otp,
			purpose: purpose,
			mockSetup: func(r *MockRedisRepo) {
				r.On("SaveOTP", mock.Anything, email, otp, purpose).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "Fail - Redis Error",
			email:   email,
			otp:     otp,
			purpose: purpose,
			mockSetup: func(r *MockRedisRepo) {
				r.On("SaveOTP", mock.Anything, email, otp, purpose).Return(xerror.Internal("Redis error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepo)
			mockRedisRepo := new(MockRedisRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRedisRepo)
			}

			authUC := usecase.NewAuthUseCase(mockUserRepo, newTestTokenMaster(), time.Second, mockRedisRepo, nil, nil)

			err := authUC.SaveOTP(context.Background(), tt.email, tt.otp, tt.purpose)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRedisRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_SendOTP(t *testing.T) {
	tests := []struct {
		name          string
		purpose       string
		mockSetup     func(*MockUserRepo, *MockRedisRepo, *MockMailer)
		expectedError bool
	}{
		{
			name:    "Success - Reset Password",
			purpose: string(constants.PurposeResetPassword),
			mockSetup: func(u *MockUserRepo, r *MockRedisRepo, m *MockMailer) {
				u.On("GetByEmail", mock.Anything, "test@gmail.com").Return(&entity.User{FullName: "Test"}, nil)
				r.On("SaveOTP", mock.Anything, "test@gmail.com", mock.Anything, constants.PurposeResetPassword).Return(nil)
				m.On("SendOTP", "test@gmail.com", "Test", mock.Anything).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepo)
			mockRedisRepo := new(MockRedisRepo)
			mockMailer := new(MockMailer)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUserRepo, mockRedisRepo, mockMailer)
			}

			authUC := usecase.NewAuthUseCase(mockUserRepo, nil, time.Second, mockRedisRepo, mockMailer, nil)
			err := authUC.SendOTP(context.Background(), "test@gmail.com", tt.purpose)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
}

func TestAuthUseCase_ResetPassword(t *testing.T) {
	email := "test@gmail.com"
	newPassword := "NewPass123@"
	userID := uuid.New()

	mockUser := &entity.User{
		ID:       userID,
		Email:    email,
		Password: "old_hashed_password",
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(*MockUserRepo)
		expectedError bool
	}{
		{
			name:     "Success - Reset Password",
			email:    email,
			password: newPassword,
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByEmail", mock.Anything, email).Return(mockUser, nil)
				u.On("UpdatePassword", mock.Anything, userID.String(), mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "Fail - User Not Found",
			email:    "ghost@gmail.com",
			password: newPassword,
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByEmail", mock.Anything, "ghost@gmail.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
		},
		{
			name:     "Fail - Database Error on Get",
			email:    email,
			password: newPassword,
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByEmail", mock.Anything, email).Return(nil, xerror.Internal("DB Error"))
			},
			expectedError: true,
		},
		{
			name:     "Fail - Database Error on UpdatePassword",
			email:    email,
			password: newPassword,
			mockSetup: func(u *MockUserRepo) {
				u.On("GetByEmail", mock.Anything, email).Return(mockUser, nil)
				u.On("UpdatePassword", mock.Anything, userID.String(), mock.Anything).Return(xerror.Internal("Update Error"))
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

			err := authUC.ResetPassword(context.Background(), tt.email, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}
