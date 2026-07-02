package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK SHIFT ASSIGNMENT REPO =====
type MockShiftAssignmentRepo struct{ mock.Mock }

func (m *MockShiftAssignmentRepo) Create(ctx context.Context, shiftAssignment *entity.ShiftAssignment) (*entity.ShiftAssignment, error) {
	args := m.Called(ctx, shiftAssignment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftAssignment), args.Error(1)
}

func (m *MockShiftAssignmentRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error) {
	args := m.Called(ctx, updateData, shiftID, shiftAssignmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftAssignment), args.Error(1)
}

func (m *MockShiftAssignmentRepo) Delete(ctx context.Context, shiftID, shiftAssignmentID string) error {
	args := m.Called(ctx, shiftID, shiftAssignmentID)
	return args.Error(0)
}

func (m *MockShiftAssignmentRepo) GetByID(ctx context.Context, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error) {
	args := m.Called(ctx, shiftID, shiftAssignmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftAssignment), args.Error(1)
}

func (m *MockShiftAssignmentRepo) GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftAssignment, error) {
	args := m.Called(ctx, shiftID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShiftAssignment), args.Error(1)
}

func (m *MockShiftAssignmentRepo) IsExistByUserAndShift(ctx context.Context, shiftID, userID string) (bool, error) {
	args := m.Called(ctx, shiftID, userID)
	return args.Bool(0), args.Error(1)
}
