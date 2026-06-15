package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK SHIFT REPO =====
type MockShiftRepo struct{ mock.Mock }

func (m *MockShiftRepo) Create(ctx context.Context, shift *entity.Shift) (*entity.Shift, error) {
	args := m.Called(ctx, shift)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Shift), args.Error(1)
}

func (m *MockShiftRepo) Update(ctx context.Context, shiftID, scheID string, updateData map[string]interface{}) (*entity.Shift, error) {
	args := m.Called(ctx, shiftID, scheID, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Shift), args.Error(1)
}

func (m *MockShiftRepo) Delete(ctx context.Context, scheID, shiftID string) error {
	args := m.Called(ctx, scheID, shiftID)
	return args.Error(0)
}

func (m *MockShiftRepo) FindByID(ctx context.Context, scheID, shiftID string) (*entity.Shift, error) {
	args := m.Called(ctx, scheID, shiftID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Shift), args.Error(1)
}

func (m *MockShiftRepo) FindAllByScheduleID(ctx context.Context, scheID string) ([]*entity.Shift, error) {
	args := m.Called(ctx, scheID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Shift), args.Error(1)
}
