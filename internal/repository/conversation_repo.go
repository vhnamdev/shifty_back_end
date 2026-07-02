package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConversationRepository interface {
	Create(ctx context.Context, conversation *entity.Conversation) (*entity.Conversation, error)
	CreateParticipants(ctx context.Context, participants []*entity.Participant) error
	FindDirectByParticipants(ctx context.Context, restaurantID, userID, targetUserID string) (*entity.Conversation, error)
	GetByID(ctx context.Context, conversationID, restaurantID string) (*entity.Conversation, error)
	GetAllByUserID(ctx context.Context, userID, restaurantID string, page, limit int) ([]*entity.Conversation, int64, error)
	IsActiveParticipant(ctx context.Context, conversationID, userID string) (bool, error)
	HideForUser(ctx context.Context, conversationID, userID string) error
	CreateMessage(ctx context.Context, message *entity.Message) (*entity.Message, error)
	GetMessages(ctx context.Context, conversationID string, page, limit int) ([]*entity.Message, int64, error)
	UpdateLastMessageAt(ctx context.Context, conversationID string, lastMessageAt time.Time) error
}

type conversationRepo struct{ db *gorm.DB }

func NewConversationRepository(db *gorm.DB) ConversationRepository { return &conversationRepo{db: db} }
func (r *conversationRepo) Create(ctx context.Context, conversation *entity.Conversation) (*entity.Conversation, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(conversation).Error; err != nil {
		return nil, err
	}
	return conversation, nil
}
func (r *conversationRepo) CreateParticipants(ctx context.Context, participants []*entity.Participant) error {
	return r.db.WithContext(ctx).Create(&participants).Error
}
func (r *conversationRepo) FindDirectByParticipants(ctx context.Context, restaurantID, userID, targetUserID string) (*entity.Conversation, error) {
	var conversation entity.Conversation
	if err := r.db.WithContext(ctx).Model(&entity.Conversation{}).
		Joins("JOIN participants p1 ON p1.conversation_id = conversations.id AND p1.author_id = ? AND p1.is_deleted = ?", userID, false).
		Joins("JOIN participants p2 ON p2.conversation_id = conversations.id AND p2.author_id = ? AND p2.is_deleted = ?", targetUserID, false).
		Preload("Participants", "is_deleted = ?", false).
		Where("conversations.restaurant_id = ? AND conversations.type = ? AND conversations.is_deleted = ?", restaurantID, entity.ConvTypeDirect, false).
		First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}
func (r *conversationRepo) GetByID(ctx context.Context, conversationID, restaurantID string) (*entity.Conversation, error) {
	var conversation entity.Conversation
	if err := r.db.WithContext(ctx).Preload("Participants", "is_deleted = ?", false).Where("id = ? AND restaurant_id = ? AND is_deleted = ?", conversationID, restaurantID, false).First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}
func (r *conversationRepo) GetAllByUserID(ctx context.Context, userID, restaurantID string, page, limit int) ([]*entity.Conversation, int64, error) {
	var conversations []*entity.Conversation
	var total int64
	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&entity.Conversation{}).Joins("JOIN participants ON participants.conversation_id = conversations.id AND participants.author_id = ? AND participants.is_deleted = ?", userID, false).Where("conversations.restaurant_id = ? AND conversations.is_deleted = ?", restaurantID, false)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("Participants", "is_deleted = ?", false).Limit(limit).Offset(offset).Order("last_message_at desc nulls last, created_at desc").Find(&conversations).Error; err != nil {
		return nil, 0, err
	}
	return conversations, total, nil
}
func (r *conversationRepo) IsActiveParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Participant{}).Where("conversation_id = ? AND author_id = ? AND is_deleted = ?", conversationID, userID, false).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *conversationRepo) HideForUser(ctx context.Context, conversationID, userID string) error {
	return r.db.WithContext(ctx).Model(&entity.Participant{}).Where("conversation_id = ? AND author_id = ? AND is_deleted = ?", conversationID, userID, false).Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error
}
func (r *conversationRepo) CreateMessage(ctx context.Context, message *entity.Message) (*entity.Message, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}
func (r *conversationRepo) GetMessages(ctx context.Context, conversationID string, page, limit int) ([]*entity.Message, int64, error) {
	var messages []*entity.Message
	var total int64
	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&entity.Message{}).Where("conversation_id = ? AND is_deleted = ?", conversationID, false)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("Sender", "is_deleted = ?", false).Limit(limit).Offset(offset).Order("created_at desc").Find(&messages).Error; err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}
func (r *conversationRepo) UpdateLastMessageAt(ctx context.Context, conversationID string, lastMessageAt time.Time) error {
	return r.db.WithContext(ctx).Model(&entity.Conversation{}).Where("id = ? AND is_deleted = ?", conversationID, false).Update("last_message_at", lastMessageAt).Error
}
