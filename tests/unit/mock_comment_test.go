package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK COMMENT REPO =====
type MockCommentRepo struct{ mock.Mock }

func (m *MockCommentRepo) Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	args := m.Called(ctx, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentRepo) Update(ctx context.Context, updateData map[string]interface{}, commentID, postID string) (*entity.Comment, error) {
	args := m.Called(ctx, updateData, commentID, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentRepo) Delete(ctx context.Context, commentID, postID string) error {
	args := m.Called(ctx, commentID, postID)
	return args.Error(0)
}

func (m *MockCommentRepo) GetByID(ctx context.Context, commentID, postID string) (*entity.Comment, error) {
	args := m.Called(ctx, commentID, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentRepo) GetAllByPostID(ctx context.Context, postID string, page, limit int) ([]*entity.Comment, int64, error) {
	args := m.Called(ctx, postID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Comment), args.Get(1).(int64), args.Error(2)
}
