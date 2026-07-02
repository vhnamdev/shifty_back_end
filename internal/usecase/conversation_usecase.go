package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ConversationUseCase interface {
	CreateDirectConversation(ctx context.Context, userID, resID, targetUserID string) (*entity.Conversation, error)
	CreateGroupConversation(ctx context.Context, userID, resID, name string, avatar *string, participantIDs []string) (*entity.Conversation, error)
	FindByID(ctx context.Context, userID, resID, conversationID string) (*entity.Conversation, error)
	FindAllByUserID(ctx context.Context, userID, resID string, page, limit int) ([]*entity.Conversation, int64, error)
	HideConversation(ctx context.Context, userID, resID, conversationID string) error
	SendMessage(ctx context.Context, userID, resID, conversationID, content string, imageURL *string) (*entity.Message, error)
	FindMessages(ctx context.Context, userID, resID, conversationID string, page, limit int) ([]*entity.Message, int64, error)
}

type conversationUseCase struct {
	conversationRepo   repository.ConversationRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewConversationUseCase(conversationRepo repository.ConversationRepository, userRestaurantRepo repository.UserRestaurantRepository) ConversationUseCase {
	return &conversationUseCase{conversationRepo: conversationRepo, userRestaurantRepo: userRestaurantRepo}
}
func (u *conversationUseCase) CreateDirectConversation(ctx context.Context, userID, resID, targetUserID string) (*entity.Conversation, error) {
	if userID == targetUserID {
		return nil, xerror.BadRequest("Can not create direct conversation with yourself")
	}
	parsedResID, err := uuid.Parse(resID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid restaurant ID")
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid user ID")
	}
	parsedTargetID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid target user ID")
	}
	if err := u.requireRestaurantMembers(ctx, resID, userID, targetUserID); err != nil {
		return nil, err
	}
	existingConversation, err := u.conversationRepo.FindDirectByParticipants(ctx, resID, userID, targetUserID)
	if err == nil {
		return existingConversation, nil
	}
	if !utils.IsRecordNotFoundError(err) {
		return nil, xerror.Internal("Database failed")
	}
	conversation, err := u.conversationRepo.Create(ctx, &entity.Conversation{Type: entity.ConvTypeDirect, RestaurantID: parsedResID})
	if err != nil {
		return nil, xerror.Internal("Can not create conversation")
	}
	participants := []*entity.Participant{{ConversationID: conversation.ID, AuthorID: parsedUserID}, {ConversationID: conversation.ID, AuthorID: parsedTargetID}}
	if err := u.conversationRepo.CreateParticipants(ctx, participants); err != nil {
		return nil, xerror.Internal("Can not create participants")
	}
	conversation.Participants = []entity.Participant{*participants[0], *participants[1]}
	return conversation, nil
}
func (u *conversationUseCase) CreateGroupConversation(ctx context.Context, userID, resID, name string, avatar *string, participantIDs []string) (*entity.Conversation, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, xerror.BadRequest("Conversation name is required")
	}
	parsedResID, err := uuid.Parse(resID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid restaurant ID")
	}
	participantMap := map[string]struct{}{userID: {}}
	for _, participantID := range participantIDs {
		participantID = strings.TrimSpace(participantID)
		if participantID != "" {
			participantMap[participantID] = struct{}{}
		}
	}
	if len(participantMap) < 3 {
		return nil, xerror.BadRequest("Group conversation must have at least 3 participants")
	}
	allParticipantIDs := make([]string, 0, len(participantMap))
	participants := make([]*entity.Participant, 0, len(participantMap))
	for participantID := range participantMap {
		parsedParticipantID, err := uuid.Parse(participantID)
		if err != nil {
			return nil, xerror.BadRequest("Invalid participant ID")
		}
		allParticipantIDs = append(allParticipantIDs, participantID)
		participants = append(participants, &entity.Participant{AuthorID: parsedParticipantID})
	}
	if err := u.requireRestaurantMembers(ctx, resID, allParticipantIDs...); err != nil {
		return nil, err
	}
	conversation, err := u.conversationRepo.Create(ctx, &entity.Conversation{Type: entity.ConvTypeGroup, Name: &name, Avatar: avatar, RestaurantID: parsedResID})
	if err != nil {
		return nil, xerror.Internal("Can not create conversation")
	}
	for _, participant := range participants {
		participant.ConversationID = conversation.ID
	}
	if err := u.conversationRepo.CreateParticipants(ctx, participants); err != nil {
		return nil, xerror.Internal("Can not create participants")
	}
	conversation.Participants = make([]entity.Participant, 0, len(participants))
	for _, participant := range participants {
		conversation.Participants = append(conversation.Participants, *participant)
	}
	return conversation, nil
}
func (u *conversationUseCase) FindByID(ctx context.Context, userID, resID, conversationID string) (*entity.Conversation, error) {
	conversation, err := u.getAccessibleConversation(ctx, userID, resID, conversationID)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}
func (u *conversationUseCase) FindAllByUserID(ctx context.Context, userID, resID string, page, limit int) ([]*entity.Conversation, int64, error) {
	if err := u.requireRestaurantMembers(ctx, resID, userID); err != nil {
		return nil, 0, err
	}
	page, limit = normalizeConversationPagination(page, limit)
	conversations, total, err := u.conversationRepo.GetAllByUserID(ctx, userID, resID, page, limit)
	if err != nil {
		return nil, 0, xerror.Internal("Database failed")
	}
	return conversations, total, nil
}
func (u *conversationUseCase) HideConversation(ctx context.Context, userID, resID, conversationID string) error {
	if _, err := u.getAccessibleConversation(ctx, userID, resID, conversationID); err != nil {
		return err
	}
	if err := u.conversationRepo.HideForUser(ctx, conversationID, userID); err != nil {
		return xerror.Internal("Can not hide conversation")
	}
	return nil
}
func (u *conversationUseCase) SendMessage(ctx context.Context, userID, resID, conversationID, content string, imageURL *string) (*entity.Message, error) {
	content = strings.TrimSpace(content)
	if imageURL != nil {
		trimmedImage := strings.TrimSpace(*imageURL)
		imageURL = &trimmedImage
	}
	if content == "" && (imageURL == nil || *imageURL == "") {
		return nil, xerror.BadRequest("Message must have content or image")
	}
	if _, err := u.getAccessibleConversation(ctx, userID, resID, conversationID); err != nil {
		return nil, err
	}
	parsedConversationID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid conversation ID")
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid user ID")
	}
	message, err := u.conversationRepo.CreateMessage(ctx, &entity.Message{ConversationID: parsedConversationID, SenderID: parsedUserID, Content: content, ImageUrl: imageURL})
	if err != nil {
		return nil, xerror.Internal("Can not create message")
	}
	if err := u.conversationRepo.UpdateLastMessageAt(ctx, conversationID, time.Now()); err != nil {
		return nil, xerror.Internal("Can not update conversation")
	}
	return message, nil
}
func (u *conversationUseCase) FindMessages(ctx context.Context, userID, resID, conversationID string, page, limit int) ([]*entity.Message, int64, error) {
	if _, err := u.getAccessibleConversation(ctx, userID, resID, conversationID); err != nil {
		return nil, 0, err
	}
	page, limit = normalizeConversationPagination(page, limit)
	messages, total, err := u.conversationRepo.GetMessages(ctx, conversationID, page, limit)
	if err != nil {
		return nil, 0, xerror.Internal("Database failed")
	}
	return messages, total, nil
}
func (u *conversationUseCase) getAccessibleConversation(ctx context.Context, userID, resID, conversationID string) (*entity.Conversation, error) {
	conversation, err := u.conversationRepo.GetByID(ctx, conversationID, resID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Conversation is not found")
		}
		return nil, xerror.Internal("Database failed")
	}
	isParticipant, err := u.conversationRepo.IsActiveParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, xerror.Internal("Can not check participant")
	}
	if !isParticipant {
		return nil, xerror.Forbidden("You are not allowed to access conversation")
	}
	return conversation, nil
}
func (u *conversationUseCase) requireRestaurantMembers(ctx context.Context, resID string, userIDs ...string) error {
	for _, userID := range userIDs {
		isMember, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)
		if err != nil {
			return xerror.Internal("Can not check member")
		}
		if !isMember {
			return xerror.Forbidden("User is not in restaurant")
		}
	}
	return nil
}
func normalizeConversationPagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return page, limit
}
