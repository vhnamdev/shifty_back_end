package unit_test

import (
	"context"
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

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
func (m *MockUserRepo) UpdateImage(ctx context.Context, id, imageURL string) (*entity.User, error) {
	args := m.Called(ctx, id, imageURL)
	return args.Get(0).(*entity.User), args.Error(1)
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
