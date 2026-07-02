package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"
)

func MapConversationEntityToModel(conversation *entity.Conversation) (*model.Conversation, error) {
	if conversation == nil {
		return nil, xerror.BadRequest("Conversation is not valid")
	}
	participants := make([]*model.Participant, 0, len(conversation.Participants))
	for _, participant := range conversation.Participants {
		mappedParticipant, err := MapParticipantEntityToModel(&participant)
		if err != nil {
			return nil, err
		}
		participants = append(participants, mappedParticipant)
	}
	return &model.Conversation{ID: conversation.ID.String(), Type: model.ConversationType(conversation.Type), Name: conversation.Name, Avatar: conversation.Avatar, RestaurantID: conversation.RestaurantID.String(), LastMessageAt: conversation.LastMessageAt, Participants: participants, IsDeleted: conversation.IsDeleted, CreatedAt: conversation.CreatedAt, UpdatedAt: conversation.UpdatedAt, DeletedAt: conversation.DeletedAt}, nil
}
func MapParticipantEntityToModel(participant *entity.Participant) (*model.Participant, error) {
	if participant == nil {
		return nil, xerror.BadRequest("Participant is not valid")
	}
	return &model.Participant{ID: participant.ID.String(), ConversationID: participant.ConversationID.String(), AuthorID: participant.AuthorID.String(), IsDeleted: participant.IsDeleted, CreatedAt: participant.CreatedAt, UpdatedAt: participant.UpdatedAt, DeletedAt: participant.DeletedAt}, nil
}
func MapMessageEntityToModel(message *entity.Message) (*model.Message, error) {
	if message == nil {
		return nil, xerror.BadRequest("Message is not valid")
	}
	return &model.Message{ID: message.ID.String(), ConversationID: message.ConversationID.String(), SenderID: message.SenderID.String(), Content: message.Content, ImageURL: message.ImageUrl, IsDeleted: message.IsDeleted, CreatedAt: message.CreatedAt, UpdatedAt: message.UpdatedAt, DeletedAt: message.DeletedAt}, nil
}
func MapConversationsToPagination(conversations []*entity.Conversation, total int64, page, limit int) (*model.ConversationPagination, error) {
	conversationModels := make([]*model.Conversation, 0, len(conversations))
	for _, conversation := range conversations {
		mappedConversation, err := MapConversationEntityToModel(conversation)
		if err != nil {
			return nil, err
		}
		conversationModels = append(conversationModels, mappedConversation)
	}
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return &model.ConversationPagination{Data: conversationModels, Total: int(total), CurrentPage: page, TotalPages: totalPages}, nil
}
func MapMessagesToPagination(messages []*entity.Message, total int64, page, limit int) (*model.MessagePagination, error) {
	messageModels := make([]*model.Message, 0, len(messages))
	for _, message := range messages {
		mappedMessage, err := MapMessageEntityToModel(message)
		if err != nil {
			return nil, err
		}
		messageModels = append(messageModels, mappedMessage)
	}
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return &model.MessagePagination{Data: messageModels, Total: int(total), CurrentPage: page, TotalPages: totalPages}, nil
}
