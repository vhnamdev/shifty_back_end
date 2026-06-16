package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// --- MOCK SCHEDULE REPO ---
type MockScheduleRepo struct {
	mock.Mock
}

func (m *MockScheduleRepo) Create(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	args := m.Called(ctx, schedule)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Schedule), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockScheduleRepo) Update(ctx context.Context, resID, scheID string, updateData map[string]interface{}) (*entity.Schedule, error) {
	args := m.Called(ctx, resID, scheID, updateData)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Schedule), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockScheduleRepo) Delete(ctx context.Context, scheID, resID string) error {
	args := m.Called(ctx, scheID, resID)
	return args.Error(0)
}

func (m *MockScheduleRepo) FindByID(ctx context.Context, scheID, resID string) (*entity.Schedule, error) {
	args := m.Called(ctx, scheID, resID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Schedule), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockScheduleRepo) FindAllByResID(ctx context.Context, resID string) ([]*entity.Schedule, error) {
	args := m.Called(ctx, resID)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Schedule), args.Error(1)
	}
	return nil, args.Error(1)
}
