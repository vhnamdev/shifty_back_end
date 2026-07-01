package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK POST REPO =====
type MockPostRepo struct{ mock.Mock }

func (m *MockPostRepo) Create(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepo) Update(ctx context.Context, updateData map[string]interface{}, postID, restaurantID string) (*entity.Post, error) {
	args := m.Called(ctx, updateData, postID, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepo) Delete(ctx context.Context, postID, restaurantID string) error {
	args := m.Called(ctx, postID, restaurantID)
	return args.Error(0)
}

func (m *MockPostRepo) GetByID(ctx context.Context, postID, restaurantID string) (*entity.Post, error) {
	args := m.Called(ctx, postID, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepo) GetAllByRestaurantID(ctx context.Context, restaurantID string, page, limit int) ([]*entity.Post, int64, error) {
	args := m.Called(ctx, restaurantID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Post), args.Get(1).(int64), args.Error(2)
}
