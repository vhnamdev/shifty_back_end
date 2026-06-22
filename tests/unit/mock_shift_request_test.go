package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK SHIFT REQUEST REPO =====
type MockShiftRequestRepo struct{ mock.Mock }

func (m *MockShiftRequestRepo) Create(ctx context.Context, shiftRequest *entity.ShiftRequest) (*entity.ShiftRequest, error) {
	args := m.Called(ctx, shiftRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRequest), args.Error(1)
}

func (m *MockShiftRequestRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftRequestID string) (*entity.ShiftRequest, error) {
	args := m.Called(ctx, updateData, shiftID, shiftRequestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRequest), args.Error(1)
}

func (m *MockShiftRequestRepo) Delete(ctx context.Context, shiftID, shiftRequestID string) error {
	args := m.Called(ctx, shiftID, shiftRequestID)
	return args.Error(0)
}

func (m *MockShiftRequestRepo) GetByID(ctx context.Context, shiftID, shiftRequestID string) (*entity.ShiftRequest, error) {
	args := m.Called(ctx, shiftID, shiftRequestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShiftRequest), args.Error(1)
}

func (m *MockShiftRequestRepo) GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftRequest, error) {
	args := m.Called(ctx, shiftID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShiftRequest), args.Error(1)
}

func (m *MockShiftRequestRepo) GetAllByUserID(ctx context.Context, userID string) ([]*entity.ShiftRequest, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShiftRequest), args.Error(1)
}
