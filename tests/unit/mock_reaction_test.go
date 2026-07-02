package unit_test

import (
	"context"
	"shifty-backend/internal/entity"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK REACTION REPO =====
type MockReactionRepo struct{ mock.Mock }

func (m *MockReactionRepo) Create(ctx context.Context, reaction *entity.Reaction) (*entity.Reaction, error) {
	args := m.Called(ctx, reaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Reaction), args.Error(1)
}

func (m *MockReactionRepo) UpdateType(ctx context.Context, reactionID, reactionType string) (*entity.Reaction, error) {
	args := m.Called(ctx, reactionID, reactionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Reaction), args.Error(1)
}

func (m *MockReactionRepo) Delete(ctx context.Context, reactionID string) error {
	args := m.Called(ctx, reactionID)
	return args.Error(0)
}

func (m *MockReactionRepo) GetByPostAndAuthor(ctx context.Context, postID, authorID string) (*entity.Reaction, error) {
	args := m.Called(ctx, postID, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Reaction), args.Error(1)
}

func (m *MockReactionRepo) CountByPostID(ctx context.Context, postID string) (map[string]int, int, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).(map[string]int), args.Int(1), args.Error(2)
}
