package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK SHIFT REQUIREMENT REPO =====
type MockShiftRequirementRepo struct{ mock.Mock }

func (m *MockShiftRequirementRepo) Create(ctx context.Context, shiftRequirement *entity.ShiftRequirement) (*entity.ShiftRequirement, error) {
	args := m.Called(ctx, shiftRequirement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRequirement), args.Error(1)
}

func (m *MockShiftRequirementRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error) {
	args := m.Called(ctx, updateData, shiftID, shiftRequireID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRequirement), args.Error(1)
}

func (m *MockShiftRequirementRepo) Delete(ctx context.Context, shiftID, shiftRequireID string) error {
	args := m.Called(ctx, shiftID, shiftRequireID)
	return args.Error(0)
}

func (m *MockShiftRequirementRepo) GetByID(ctx context.Context, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error) {
	args := m.Called(ctx, shiftID, shiftRequireID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRequirement), args.Error(1)
}

func (m *MockShiftRequirementRepo) GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftRequirement, error) {
	args := m.Called(ctx, shiftID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShiftRequirement), args.Error(1)
}
