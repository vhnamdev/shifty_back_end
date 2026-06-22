package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK SHIFT RULE REPO =====
type MockShiftRuleRepo struct{ mock.Mock }

func (m *MockShiftRuleRepo) Create(ctx context.Context, shiftRule *entity.ShiftRule) (*entity.ShiftRule, error) {
	args := m.Called(ctx, shiftRule)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRule), args.Error(1)
}

func (m *MockShiftRuleRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftRuleID, restaurantID string) (*entity.ShiftRule, error) {
	args := m.Called(ctx, updateData, shiftRuleID, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRule), args.Error(1)
}

func (m *MockShiftRuleRepo) Delete(ctx context.Context, shiftRuleID, restaurantID string) error {
	args := m.Called(ctx, shiftRuleID, restaurantID)
	return args.Error(0)
}

func (m *MockShiftRuleRepo) GetByID(ctx context.Context, shiftRuleID, restaurantID string) (*entity.ShiftRule, error) {
	args := m.Called(ctx, shiftRuleID, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRule), args.Error(1)
}

func (m *MockShiftRuleRepo) GetAllByRestaurantID(ctx context.Context, restaurantID string) ([]*entity.ShiftRule, error) {
	args := m.Called(ctx, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShiftRule), args.Error(1)
}
