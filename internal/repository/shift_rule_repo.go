package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShiftRuleRepository interface {
	Create(ctx context.Context, shiftRule *entity.ShiftRule) (*entity.ShiftRule, error)
	Update(ctx context.Context, updateData map[string]interface{}, shiftRuleID, restaurantID string) (*entity.ShiftRule, error)
	Delete(ctx context.Context, shiftRuleID, restaurantID string) error
	GetByID(ctx context.Context, shiftRuleID, restaurantID string) (*entity.ShiftRule, error)
	GetAllByRestaurantID(ctx context.Context, restaurantID string) ([]*entity.ShiftRule, error)
}

type shiftRuleRepo struct {
	db *gorm.DB
}

func NewShiftRuleRepository(db *gorm.DB) ShiftRuleRepository {
	return &shiftRuleRepo{
		db: db,
	}
}

func (r *shiftRuleRepo) Create(ctx context.Context, shiftRule *entity.ShiftRule) (*entity.ShiftRule, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(shiftRule).Error; err != nil {
		return nil, err
	}
	return shiftRule, nil
}

func (r *shiftRuleRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftRuleID, restaurantID string) (*entity.ShiftRule, error) {
	var shiftRule entity.ShiftRule

	if err := r.db.
		WithContext(ctx).
		Model(&shiftRule).
		Clauses(clause.Returning{}).
		Where("id = ? AND restaurant_id = ?", shiftRuleID, restaurantID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &shiftRule, nil
}

func (r *shiftRuleRepo) Delete(ctx context.Context, shiftRuleID, restaurantID string) error {
	if err := r.db.Model(&entity.ShiftRule{}).Where("id = ? AND restaurant_id = ?", shiftRuleID, restaurantID).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *shiftRuleRepo) GetByID(ctx context.Context, shiftRuleID, restaurantID string) (*entity.ShiftRule, error) {
	var shiftRule entity.ShiftRule

	if err := r.db.
		WithContext(ctx).
		Preload("Restaurant", "is_deleted = ?", false).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", shiftRuleID, restaurantID, false).
		First(&shiftRule).Error; err != nil {
		return nil, err
	}
	return &shiftRule, nil
}

func (r *shiftRuleRepo) GetAllByRestaurantID(ctx context.Context, restaurantID string) ([]*entity.ShiftRule, error) {
	var shiftRules []*entity.ShiftRule

	if err := r.db.WithContext(ctx).Preload("Restaurant", "is_deleted = ?", false).Where("restaurant_id = ? AND is_deleted = ?", restaurantID, false).Find(&shiftRules).Error; err != nil {
		return nil, err
	}
	return shiftRules, nil
}
