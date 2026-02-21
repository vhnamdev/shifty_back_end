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

// Helper tạo Token dùng chung
func newTestTokenMaster() *token.TokenMaster {
	return token.NewToken(
		"secret_access_test_key",
		"secret_refresh_test_key",
		time.Minute*15,
		time.Hour*24,
	)
}

// ===== MOCK USER REPO =====
type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) IsEmailExist(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
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
	return m.Called(ctx, user).Error(0)
}
func (m *MockUserRepo) Update(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *MockUserRepo) UpdatePassword(ctx context.Context, id, newPassword string) error {
	return m.Called(ctx, id, newPassword).Error(0)
}
func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockUserRepo) GetRestaurantMembers(ctx context.Context, page int, limit int, restaurantID string, filter *dto.UserFilter) ([]*entity.User, int64, error) {
	args := m.Called(ctx, page, limit, restaurantID, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.User), int64(args.Int(1)), args.Error(2)
}

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

// ===== MOCK MAILER =====
type MockMailer struct{ mock.Mock }

func (m *MockMailer) SendOTP(email, name, otp string) error {
	return m.Called(email, name, otp).Error(0)
}

// ===== MOCK USER RESTAURANT REPO =====
type MockUserRestaurantRepo struct{ mock.Mock }

func (m *MockUserRestaurantRepo) Create(ctx context.Context, userRes *entity.UserRestaurant) (*entity.UserRestaurant, error) {
	args := m.Called(ctx, userRes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserRestaurant), args.Error(1)
}
func (m *MockUserRestaurantRepo) CheckUserInRestaurant(ctx context.Context, userID, resID string) (bool, error) {
	args := m.Called(ctx, userID, resID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRestaurantRepo) CheckAuthority(ctx context.Context, targetID, requestID, resID string) (bool, error) {
	args := m.Called(ctx, targetID, requestID, resID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRestaurantRepo) CheckAuthorityToUpdate(ctx context.Context, userID, resID string) (bool, error) {
	args := m.Called(ctx, userID, resID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRestaurantRepo) CheckAuthorityToDelete(ctx context.Context, userID, resID string) (bool, error) {
	args := m.Called(ctx, userID, resID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRestaurantRepo) Update(ctx context.Context, userID, resID string, updateData map[string]interface{}) (*entity.UserRestaurant, error) {
	args := m.Called(ctx, userID, resID, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserRestaurant), args.Error(1)
}

// ===== MOCK RESTAURANT REPO =====
type MockRestaurantRepo struct{ mock.Mock }

func (m *MockRestaurantRepo) Create(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {
	args := m.Called(ctx, restaurant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Restaurant), args.Error(1)
}
func (m *MockRestaurantRepo) Update(ctx context.Context, resID string, updateData map[string]interface{}) (*entity.Restaurant, error) {
	args := m.Called(ctx, resID, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Restaurant), args.Error(1)
}
func (m *MockRestaurantRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockRestaurantRepo) GetByID(ctx context.Context, id string) (*entity.Restaurant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Restaurant), args.Error(1)
}
func (m *MockRestaurantRepo) GetMyRestaurants(ctx context.Context, userID string) ([]*entity.Restaurant, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Restaurant), args.Error(1)
}

// ===== MOCK POSITION REPO =====
type MockPositionRepo struct{ mock.Mock }

func (m *MockPositionRepo) Create(ctx context.Context, position *entity.Position) (*entity.Position, error) {
	args := m.Called(ctx, position)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Position), args.Error(1)
}

// ===== MOCK TRANSACTOR =====
type MockTransactor struct{ mock.Mock }

func (m *MockTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
