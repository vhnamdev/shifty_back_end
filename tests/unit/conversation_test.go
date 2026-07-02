package unit_test

import (
	"context"
	"testing"
	"time"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupConversationUseCase() (*MockConversationRepo, *MockUserRestaurantRepo, usecase.ConversationUseCase) {
	mockConversationRepo := new(MockConversationRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)
	u := usecase.NewConversationUseCase(mockConversationRepo, mockUserResRepo)
	return mockConversationRepo, mockUserResRepo, u
}

func TestConversationUseCase_CreateDirectConversation(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	targetUserID := uuid.New()
	resID := uuid.New()
	t.Run("Success Reuse Existing", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		existingConversation := &entity.Conversation{ID: uuid.New(), Type: entity.ConvTypeDirect, RestaurantID: resID}
		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockUserResRepo.On("CheckUserInRestaurant", ctx, targetUserID.String(), resID.String()).Return(true, nil)
		mockConversationRepo.On("FindDirectByParticipants", ctx, resID.String(), userID.String(), targetUserID.String()).Return(existingConversation, nil)
		res, err := u.CreateDirectConversation(ctx, userID.String(), resID.String(), targetUserID.String())
		assert.NoError(t, err)
		assert.Equal(t, existingConversation, res)
		mockConversationRepo.AssertNotCalled(t, "Create")
		mockConversationRepo.AssertNotCalled(t, "CreateParticipants")
		mockUserResRepo.AssertExpectations(t)
		mockConversationRepo.AssertExpectations(t)
	})
	t.Run("Success Create New", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		createdConversation := &entity.Conversation{ID: uuid.New(), Type: entity.ConvTypeDirect, RestaurantID: resID}
		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockUserResRepo.On("CheckUserInRestaurant", ctx, targetUserID.String(), resID.String()).Return(true, nil)
		mockConversationRepo.On("FindDirectByParticipants", ctx, resID.String(), userID.String(), targetUserID.String()).Return(nil, gorm.ErrRecordNotFound)
		mockConversationRepo.On("Create", ctx, mock.MatchedBy(func(conversation *entity.Conversation) bool {
			return conversation.Type == entity.ConvTypeDirect && conversation.RestaurantID == resID
		})).Return(createdConversation, nil)
		mockConversationRepo.On("CreateParticipants", ctx, mock.MatchedBy(func(participants []*entity.Participant) bool {
			if len(participants) != 2 {
				return false
			}
			ids := map[uuid.UUID]bool{participants[0].AuthorID: true, participants[1].AuthorID: true}
			return ids[userID] && ids[targetUserID]
		})).Return(nil)
		res, err := u.CreateDirectConversation(ctx, userID.String(), resID.String(), targetUserID.String())
		assert.NoError(t, err)
		assert.Equal(t, createdConversation.ID, res.ID)
		assert.Len(t, res.Participants, 2)
		mockUserResRepo.AssertExpectations(t)
		mockConversationRepo.AssertExpectations(t)
	})
	t.Run("Fail Self Conversation", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		res, err := u.CreateDirectConversation(ctx, userID.String(), resID.String(), userID.String())
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Can not create direct conversation with yourself")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockConversationRepo.AssertNotCalled(t, "FindDirectByParticipants")
	})
}

func TestConversationUseCase_CreateGroupConversation(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	memberID := uuid.New()
	otherMemberID := uuid.New()
	resID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		createdConversation := &entity.Conversation{ID: uuid.New(), Type: entity.ConvTypeGroup, RestaurantID: resID}
		mockUserResRepo.On("CheckUserInRestaurant", ctx, mock.AnythingOfType("string"), resID.String()).Return(true, nil).Times(3)
		mockConversationRepo.On("Create", ctx, mock.MatchedBy(func(conversation *entity.Conversation) bool {
			return conversation.Type == entity.ConvTypeGroup && conversation.Name != nil && *conversation.Name == "Ops" && conversation.RestaurantID == resID
		})).Return(createdConversation, nil)
		mockConversationRepo.On("CreateParticipants", ctx, mock.MatchedBy(func(participants []*entity.Participant) bool { return len(participants) == 3 })).Return(nil)
		res, err := u.CreateGroupConversation(ctx, userID.String(), resID.String(), " Ops ", nil, []string{memberID.String(), otherMemberID.String()})
		assert.NoError(t, err)
		assert.Equal(t, createdConversation.ID, res.ID)
		assert.Len(t, res.Participants, 3)
		mockUserResRepo.AssertExpectations(t)
		mockConversationRepo.AssertExpectations(t)
	})
	t.Run("Fail Too Few Participants", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		res, err := u.CreateGroupConversation(ctx, userID.String(), resID.String(), "Ops", nil, []string{memberID.String()})
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Group conversation must have at least 3 participants")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockConversationRepo.AssertNotCalled(t, "Create")
	})
	t.Run("Fail Non Restaurant Participant", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		mockUserResRepo.On("CheckUserInRestaurant", ctx, mock.AnythingOfType("string"), resID.String()).Return(true, nil).Twice()
		mockUserResRepo.On("CheckUserInRestaurant", ctx, mock.AnythingOfType("string"), resID.String()).Return(false, nil).Once()
		res, err := u.CreateGroupConversation(ctx, userID.String(), resID.String(), "Ops", nil, []string{memberID.String(), otherMemberID.String()})
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "User is not in restaurant")
		mockConversationRepo.AssertNotCalled(t, "Create")
	})
}

func TestConversationUseCase_SendMessage(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	resID := uuid.New()
	conversationID := uuid.New()
	conversation := &entity.Conversation{ID: conversationID, RestaurantID: resID, Type: entity.ConvTypeDirect}
	t.Run("Success", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		expectedMessage := &entity.Message{ID: uuid.New(), ConversationID: conversationID, SenderID: userID, Content: "Hello"}
		mockConversationRepo.On("GetByID", ctx, conversationID.String(), resID.String()).Return(conversation, nil)
		mockConversationRepo.On("IsActiveParticipant", ctx, conversationID.String(), userID.String()).Return(true, nil)
		mockConversationRepo.On("CreateMessage", ctx, mock.MatchedBy(func(message *entity.Message) bool {
			return message.ConversationID == conversationID && message.SenderID == userID && message.Content == "Hello"
		})).Return(expectedMessage, nil)
		mockConversationRepo.On("UpdateLastMessageAt", ctx, conversationID.String(), mock.AnythingOfType("time.Time")).Return(nil)
		res, err := u.SendMessage(ctx, userID.String(), resID.String(), conversationID.String(), " Hello ", nil)
		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, res)
		mockUserResRepo.AssertExpectations(t)
		mockConversationRepo.AssertExpectations(t)
	})
	t.Run("Fail Empty Message", func(t *testing.T) {
		mockConversationRepo, mockUserResRepo, u := setupConversationUseCase()
		res, err := u.SendMessage(ctx, userID.String(), resID.String(), conversationID.String(), "   ", nil)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Message must have content or image")
		mockConversationRepo.AssertNotCalled(t, "GetByID")
		mockUserResRepo.AssertExpectations(t)
	})
	t.Run("Fail Non Participant", func(t *testing.T) {
		mockConversationRepo, _, u := setupConversationUseCase()
		mockConversationRepo.On("GetByID", ctx, conversationID.String(), resID.String()).Return(conversation, nil)
		mockConversationRepo.On("IsActiveParticipant", ctx, conversationID.String(), userID.String()).Return(false, nil)
		res, err := u.SendMessage(ctx, userID.String(), resID.String(), conversationID.String(), "Hello", nil)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to access conversation")
		mockConversationRepo.AssertNotCalled(t, "CreateMessage")
	})
}

func TestConversationUseCase_FindAndHide(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	resID := uuid.New()
	conversationID := uuid.New()
	conversation := &entity.Conversation{ID: conversationID, RestaurantID: resID, Type: entity.ConvTypeGroup, LastMessageAt: ptrTime(time.Now())}
	t.Run("Find Messages Success", func(t *testing.T) {
		mockConversationRepo, _, u := setupConversationUseCase()
		messages := []*entity.Message{{ID: uuid.New(), ConversationID: conversationID, SenderID: userID, Content: "Hello"}}
		mockConversationRepo.On("GetByID", ctx, conversationID.String(), resID.String()).Return(conversation, nil)
		mockConversationRepo.On("IsActiveParticipant", ctx, conversationID.String(), userID.String()).Return(true, nil)
		mockConversationRepo.On("GetMessages", ctx, conversationID.String(), 1, 10).Return(messages, int64(1), nil)
		res, total, err := u.FindMessages(ctx, userID.String(), resID.String(), conversationID.String(), 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, messages, res)
		mockConversationRepo.AssertExpectations(t)
	})
	t.Run("Hide Conversation For Current User", func(t *testing.T) {
		mockConversationRepo, _, u := setupConversationUseCase()
		mockConversationRepo.On("GetByID", ctx, conversationID.String(), resID.String()).Return(conversation, nil)
		mockConversationRepo.On("IsActiveParticipant", ctx, conversationID.String(), userID.String()).Return(true, nil)
		mockConversationRepo.On("HideForUser", ctx, conversationID.String(), userID.String()).Return(nil)
		err := u.HideConversation(ctx, userID.String(), resID.String(), conversationID.String())
		assert.NoError(t, err)
		mockConversationRepo.AssertExpectations(t)
	})
}

func ptrTime(value time.Time) *time.Time { return &value }
