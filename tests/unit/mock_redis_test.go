package unit_test

import (
	"context"
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/constants"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK REDIS REPO =====
type MockRedisRepo struct{ mock.Mock }

// 1. Session
func (m *MockRedisRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	return m.Called(ctx, session).Error(0)
}
func (m *MockRedisRepo) GetSession(ctx context.Context, refreshToken string) (*entity.Session, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}
func (m *MockRedisRepo) DeleteSession(ctx context.Context, userID string, refreshToken string) error {
	return m.Called(ctx, userID, refreshToken).Error(0)
}
func (m *MockRedisRepo) DeleteAllSessions(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

// 2. User Cache
func (m *MockRedisRepo) SaveUserCache(ctx context.Context, user *entity.UserCache) error {
	return m.Called(ctx, user).Error(0)
}
func (m *MockRedisRepo) GetUserCache(ctx context.Context, userID string) (*entity.UserCache, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserCache), args.Error(1)
}
func (m *MockRedisRepo) DeleteUserCache(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

// 3. OTP
func (m *MockRedisRepo) SaveOTP(ctx context.Context, email string, otp string, purpose constants.OTPPurpose) error {
	return m.Called(ctx, email, otp, purpose).Error(0)
}
func (m *MockRedisRepo) VerifyOTP(ctx context.Context, email string, inputOTP string, purpose constants.OTPPurpose) error {
	return m.Called(ctx, email, inputOTP, purpose).Error(0)
}

// 4. Online Status
func (m *MockRedisRepo) SetUserStatus(ctx context.Context, userID string, isOnline bool) error {
	return m.Called(ctx, userID, isOnline).Error(0)
}
func (m *MockRedisRepo) GetUserStatus(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisRepo) SaveInviteCode(ctx context.Context, inviteCode, email, positonID, resID string) error {
	args := m.Called(ctx, inviteCode, email, positonID, resID)
	return args.Error(0)
}

func (m *MockRedisRepo) VerifyInviteCode(ctx context.Context, email, inviteCode string) (*dto.InviteData, error) {
	args := m.Called(ctx, email, inviteCode)
	return args.Get(0).(*dto.InviteData), args.Error(1)
}
