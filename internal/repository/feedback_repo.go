package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error)
	Update(ctx context.Context, feedbackID, restaurantID string, updateData map[string]interface{}) (*entity.Feedback, error)
	Delete(ctx context.Context, feedbackID, restaurantID string) error
	GetByID(ctx context.Context, feedbackID, restaurantID string) (*entity.Feedback, error)
	GetAllByRestaurantID(ctx context.Context, restaurantID string, page, limit int) ([]*entity.Feedback, int64, error)
	GetAllByMemberID(ctx context.Context, memberID, restaurantID string, page, limit int) ([]*entity.Feedback, int64, error)
}

type feedbackRepo struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepo{db: db}
}

func (r *feedbackRepo) Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(feedback).Error; err != nil {
		return nil, err
	}
	return feedback, nil
}

func (r *feedbackRepo) Update(ctx context.Context, feedbackID, restaurantID string, updateData map[string]interface{}) (*entity.Feedback, error) {
	var feedback entity.Feedback
	if err := r.db.WithContext(ctx).
		Model(&feedback).
		Clauses(clause.Returning{}).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", feedbackID, restaurantID, false).
		Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *feedbackRepo) Delete(ctx context.Context, feedbackID, restaurantID string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Feedback{}).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", feedbackID, restaurantID, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *feedbackRepo) GetByID(ctx context.Context, feedbackID, restaurantID string) (*entity.Feedback, error) {
	var feedback entity.Feedback
	if err := r.db.WithContext(ctx).
		Preload("Member", "is_deleted = ?", false).
		Preload("Reviewer", "is_deleted = ?", false).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", feedbackID, restaurantID, false).
		First(&feedback).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *feedbackRepo) GetAllByRestaurantID(ctx context.Context, restaurantID string, page, limit int) ([]*entity.Feedback, int64, error) {
	var feedbacks []*entity.Feedback
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).
		Model(&entity.Feedback{}).
		Where("restaurant_id = ? AND is_deleted = ?", restaurantID, false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Member", "is_deleted = ?", false).
		Preload("Reviewer", "is_deleted = ?", false).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

func (r *feedbackRepo) GetAllByMemberID(ctx context.Context, memberID, restaurantID string, page, limit int) ([]*entity.Feedback, int64, error) {
	var feedbacks []*entity.Feedback
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).
		Model(&entity.Feedback{}).
		Where("member_id = ? AND restaurant_id = ? AND is_deleted = ?", memberID, restaurantID, false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Member", "is_deleted = ?", false).
		Preload("Reviewer", "is_deleted = ?", false).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}
