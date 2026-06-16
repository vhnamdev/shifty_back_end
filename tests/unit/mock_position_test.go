package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK POSITION REPO =====
type MockPositionRepo struct{ mock.Mock }

func (m *MockPositionRepo) Create(ctx context.Context, position *entity.Position) (*entity.Position, error) {
	args := m.Called(ctx, position)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Position), args.Error(1)
}

func (m *MockPositionRepo) FindByID(ctx context.Context, posID, resID string) (*entity.Position, error) {
	args := m.Called(ctx, posID, resID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Position), args.Error(1)
}

func (m *MockPositionRepo) GetAllByRestaurantID(ctx context.Context, resID string) ([]*entity.Position, error) {
	args := m.Called(ctx, resID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Position), args.Error(1)
}

func (m *MockPositionRepo) Delete(ctx context.Context, posID, resID string) error {
	args := m.Called(ctx, posID, resID)
	return args.Error(0)
}

func (m *MockPositionRepo) DeleteAllByRestaurantID(ctx context.Context, resID string) error {
	args := m.Called(ctx, resID)
	return args.Error(0)
}

func (m *MockPositionRepo) Update(ctx context.Context, posID, resID string, updateData map[string]interface{}) (*entity.Position, error) {
	args := m.Called(ctx, posID, resID, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Position), args.Error(1)
}
