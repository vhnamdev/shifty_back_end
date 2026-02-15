package unit_test

import (
	"context"
	"time"

	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/token"

	"github.com/stretchr/testify/mock"
)

func newTestTokenMaster() *token.TokenMaster {
	return token.NewToken(
		"secret_access_test_key",
		"secret_refresh_test_key",
		time.Minute*15,
		time.Hour*24,
	)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserRepo) UpdatePassword(ctx context.Context, id, newPassword string) error {
	args := m.Called(ctx, id, newPassword)
	return args.Error(0)
}
func (m *MockUserRepo) Delete(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByGoogleID(ctx context.Context, id string) (*entity.User, error) {
	return nil, nil
}

func (m *MockUserRepo) GetRestaurantMembers(ctx context.Context, page int, limit int, restaurantID string, filter *dto.UserFilter) ([]entity.User, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.User), int64(args.Int(1)), args.Error(2)
}

type MockRedisRepo struct {
	mock.Mock
}

func (m *MockRedisRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRedisRepo) SaveOTP(ctx context.Context, email, otp string, purpose constants.OTPPurpose) error {
	args := m.Called(ctx, email, otp, purpose)
	return args.Error(0)
}

func (m *MockRedisRepo) SaveUserCache(ctx context.Context, userCache *entity.UserCache) error {
	args := m.Called(ctx, userCache)
	return args.Error(0)
}

func (m *MockRedisRepo) DeleteSession(ctx context.Context, userID string, refreshToken string) error {
	args := m.Called(ctx, userID, refreshToken)
	return args.Error(0)
}

func (m *MockRedisRepo) DeleteAllSessions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRedisRepo) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRedisRepo) DeleteUserCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRedisRepo) GetSession(ctx context.Context, refreshToken string) (*entity.Session, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockRedisRepo) GetUserCache(ctx context.Context, userID string) (*entity.UserCache, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entity.UserCache), args.Error(0)
}

func (m *MockRedisRepo) SetUserStatus(ctx context.Context, userID string, isOnline bool) error {
	args := m.Called(ctx, userID, isOnline)
	return args.Error(1)
}

func (m *MockRedisRepo) GetUserStatus(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockRedisRepo) VerifyOTP(ctx context.Context, email string, inputOTP string, purpose constants.OTPPurpose) error {
	args := m.Called(ctx, email, inputOTP, purpose)
	return args.Error(0)
}

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendOTP(email, name, otp string) error {
	args := m.Called(email, name, otp)
	return args.Error(0)
}
