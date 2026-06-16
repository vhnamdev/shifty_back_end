package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

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

func (m *MockRestaurantRepo) UpdateImage(ctx context.Context, resID, imageURl string) (*entity.Restaurant, error) {
	args := m.Called(ctx, resID, imageURl)

	return args.Get(0).(*entity.Restaurant), args.Error(1)
}
