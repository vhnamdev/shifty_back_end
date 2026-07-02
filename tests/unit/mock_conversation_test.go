package unit_test

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK CONVERSATION REPO =====
type MockConversationRepo struct{ mock.Mock }

func (m *MockConversationRepo) Create(ctx context.Context, conversation *entity.Conversation) (*entity.Conversation, error) {
	args := m.Called(ctx, conversation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Conversation), args.Error(1)
}
func (m *MockConversationRepo) CreateParticipants(ctx context.Context, participants []*entity.Participant) error {
	args := m.Called(ctx, participants)
	return args.Error(0)
}
func (m *MockConversationRepo) FindDirectByParticipants(ctx context.Context, restaurantID, userID, targetUserID string) (*entity.Conversation, error) {
	args := m.Called(ctx, restaurantID, userID, targetUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Conversation), args.Error(1)
}
func (m *MockConversationRepo) GetByID(ctx context.Context, conversationID, restaurantID string) (*entity.Conversation, error) {
	args := m.Called(ctx, conversationID, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Conversation), args.Error(1)
}
func (m *MockConversationRepo) GetAllByUserID(ctx context.Context, userID, restaurantID string, page, limit int) ([]*entity.Conversation, int64, error) {
	args := m.Called(ctx, userID, restaurantID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Conversation), args.Get(1).(int64), args.Error(2)
}
func (m *MockConversationRepo) IsActiveParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	args := m.Called(ctx, conversationID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *MockConversationRepo) HideForUser(ctx context.Context, conversationID, userID string) error {
	args := m.Called(ctx, conversationID, userID)
	return args.Error(0)
}
func (m *MockConversationRepo) CreateMessage(ctx context.Context, message *entity.Message) (*entity.Message, error) {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Message), args.Error(1)
}
func (m *MockConversationRepo) GetMessages(ctx context.Context, conversationID string, page, limit int) ([]*entity.Message, int64, error) {
	args := m.Called(ctx, conversationID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Message), args.Get(1).(int64), args.Error(2)
}
func (m *MockConversationRepo) UpdateLastMessageAt(ctx context.Context, conversationID string, lastMessageAt time.Time) error {
	args := m.Called(ctx, conversationID, lastMessageAt)
	return args.Error(0)
}
