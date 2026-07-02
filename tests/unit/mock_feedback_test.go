package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK FEEDBACK REPO =====
type MockFeedbackRepo struct{ mock.Mock }

func (m *MockFeedbackRepo) Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error) {
	args := m.Called(ctx, feedback)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Feedback), args.Error(1)
}
func (m *MockFeedbackRepo) Update(ctx context.Context, feedbackID, restaurantID string, updateData map[string]interface{}) (*entity.Feedback, error) {
	args := m.Called(ctx, feedbackID, restaurantID, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Feedback), args.Error(1)
}
func (m *MockFeedbackRepo) Delete(ctx context.Context, feedbackID, restaurantID string) error {
	args := m.Called(ctx, feedbackID, restaurantID)
	return args.Error(0)
}
func (m *MockFeedbackRepo) GetByID(ctx context.Context, feedbackID, restaurantID string) (*entity.Feedback, error) {
	args := m.Called(ctx, feedbackID, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Feedback), args.Error(1)
}
func (m *MockFeedbackRepo) GetAllByRestaurantID(ctx context.Context, restaurantID string, page, limit int) ([]*entity.Feedback, int64, error) {
	args := m.Called(ctx, restaurantID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Feedback), args.Get(1).(int64), args.Error(2)
}
func (m *MockFeedbackRepo) GetAllByMemberID(ctx context.Context, memberID, restaurantID string, page, limit int) ([]*entity.Feedback, int64, error) {
	args := m.Called(ctx, memberID, restaurantID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Feedback), args.Get(1).(int64), args.Error(2)
}
