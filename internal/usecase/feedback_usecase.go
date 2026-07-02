package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"strings"

	"github.com/google/uuid"
)

type FeedbackUseCase interface {
	Create(ctx context.Context, reviewerID, resID, memberID, content string) (*entity.Feedback, error)
	Update(ctx context.Context, reviewerID, resID, feedbackID, content string) (*entity.Feedback, error)
	Delete(ctx context.Context, reviewerID, resID, feedbackID string) error
	FindByID(ctx context.Context, userID, resID, feedbackID string) (*entity.Feedback, error)
	FindAllByRestaurantID(ctx context.Context, userID, resID string, page, limit int) ([]*entity.Feedback, int64, error)
	FindAllByMemberID(ctx context.Context, userID, resID, memberID string, page, limit int) ([]*entity.Feedback, int64, error)
}

type feedbackUseCase struct {
	feedbackRepo       repository.FeedbackRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewFeedbackUseCase(feedbackRepo repository.FeedbackRepository, userRestaurantRepo repository.UserRestaurantRepository) FeedbackUseCase {
	return &feedbackUseCase{feedbackRepo: feedbackRepo, userRestaurantRepo: userRestaurantRepo}
}

func (u *feedbackUseCase) Create(ctx context.Context, reviewerID, resID, memberID, content string) (*entity.Feedback, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, xerror.BadRequest("Feedback content is required")
	}

	parsedReviewerID, err := uuid.Parse(reviewerID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid reviewer ID")
	}
	parsedResID, err := uuid.Parse(resID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid restaurant ID")
	}
	parsedMemberID, err := uuid.Parse(memberID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid member ID")
	}

	if err := u.requireManagement(ctx, reviewerID, resID, "create feedback"); err != nil {
		return nil, err
	}

	isMember, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, memberID, resID)
	if err != nil {
		return nil, xerror.Internal("Can not check member")
	}
	if !isMember {
		return nil, xerror.Forbidden("Member is not in restaurant")
	}

	feedback := &entity.Feedback{Content: content, RestaurantID: parsedResID, MemberID: parsedMemberID, ReviewerID: parsedReviewerID}
	createdFeedback, err := u.feedbackRepo.Create(ctx, feedback)
	if err != nil {
		return nil, xerror.Internal("Can not create feedback")
	}
	return createdFeedback, nil
}

func (u *feedbackUseCase) Update(ctx context.Context, reviewerID, resID, feedbackID, content string) (*entity.Feedback, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, xerror.BadRequest("Feedback content is required")
	}
	if err := u.requireManagement(ctx, reviewerID, resID, "update feedback"); err != nil {
		return nil, err
	}
	if _, err := u.getFeedback(ctx, feedbackID, resID); err != nil {
		return nil, err
	}
	updatedFeedback, err := u.feedbackRepo.Update(ctx, feedbackID, resID, map[string]interface{}{"content": content})
	if err != nil {
		return nil, xerror.Internal("Can not update feedback")
	}
	return updatedFeedback, nil
}

func (u *feedbackUseCase) Delete(ctx context.Context, reviewerID, resID, feedbackID string) error {
	if err := u.requireManagement(ctx, reviewerID, resID, "delete feedback"); err != nil {
		return err
	}
	if _, err := u.getFeedback(ctx, feedbackID, resID); err != nil {
		return err
	}
	if err := u.feedbackRepo.Delete(ctx, feedbackID, resID); err != nil {
		return xerror.Internal("Can not delete feedback")
	}
	return nil
}

func (u *feedbackUseCase) FindByID(ctx context.Context, userID, resID, feedbackID string) (*entity.Feedback, error) {
	feedback, err := u.getFeedback(ctx, feedbackID, resID)
	if err != nil {
		return nil, err
	}
	if err := u.canViewFeedback(ctx, userID, resID, feedback.MemberID.String()); err != nil {
		return nil, err
	}
	return feedback, nil
}

func (u *feedbackUseCase) FindAllByRestaurantID(ctx context.Context, userID, resID string, page, limit int) ([]*entity.Feedback, int64, error) {
	if err := u.requireManagement(ctx, userID, resID, "view feedbacks"); err != nil {
		return nil, 0, err
	}
	page, limit = normalizePagination(page, limit)
	feedbacks, total, err := u.feedbackRepo.GetAllByRestaurantID(ctx, resID, page, limit)
	if err != nil {
		return nil, 0, xerror.Internal("Database failed")
	}
	return feedbacks, total, nil
}

func (u *feedbackUseCase) FindAllByMemberID(ctx context.Context, userID, resID, memberID string, page, limit int) ([]*entity.Feedback, int64, error) {
	if err := u.canViewFeedback(ctx, userID, resID, memberID); err != nil {
		return nil, 0, err
	}
	isMember, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, memberID, resID)
	if err != nil {
		return nil, 0, xerror.Internal("Can not check member")
	}
	if !isMember {
		return nil, 0, xerror.Forbidden("Member is not in restaurant")
	}
	page, limit = normalizePagination(page, limit)
	feedbacks, total, err := u.feedbackRepo.GetAllByMemberID(ctx, memberID, resID, page, limit)
	if err != nil {
		return nil, 0, xerror.Internal("Database failed")
	}
	return feedbacks, total, nil
}

func (u *feedbackUseCase) getFeedback(ctx context.Context, feedbackID, resID string) (*entity.Feedback, error) {
	feedback, err := u.feedbackRepo.GetByID(ctx, feedbackID, resID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Feedback is not found")
		}
		return nil, xerror.Internal("Database failed")
	}
	return feedback, nil
}

func (u *feedbackUseCase) requireManagement(ctx context.Context, userID, resID, action string) error {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)
	if err != nil {
		return xerror.Internal("Can not check authority")
	}
	if !isAuthority {
		return xerror.Forbidden("You are not allowed to " + action)
	}
	return nil
}

func (u *feedbackUseCase) canViewFeedback(ctx context.Context, userID, resID, memberID string) error {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)
	if err != nil {
		return xerror.Internal("Can not check authority")
	}
	if isAuthority {
		return nil
	}
	isMember, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)
	if err != nil {
		return xerror.Internal("Can not check authority")
	}
	if !isMember || userID != memberID {
		return xerror.Forbidden("You are not allowed to view feedback")
	}
	return nil
}

func normalizePagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return page, limit
}
