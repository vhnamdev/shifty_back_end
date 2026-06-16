package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

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
func (m *MockUserRestaurantRepo) HasManagementAuthority(ctx context.Context, userID, resID string) (bool, error) {
	args := m.Called(ctx, userID, resID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRestaurantRepo) DeleteAllByRestaurantID(ctx context.Context, resID string) error {
	args := m.Called(ctx, resID)
	return args.Error(0)
}
func (m *MockUserRestaurantRepo) DeleteAllByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockUserRestaurantRepo) SetPositionNull(ctx context.Context, posID, resID string) error {
	args := m.Called(ctx, posID, resID)
	return args.Error(0)
}
