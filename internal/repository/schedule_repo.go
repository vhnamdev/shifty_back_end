package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	Update(ctx context.Context, resID, scheID string, updateData map[string]interface{}) (*entity.Schedule, error)
	FindByID(ctx context.Context, scheID, resID string) (*entity.Schedule, error)
	FindAllByResID(ctx context.Context, resID string) ([]*entity.Schedule, error)
	Delete(ctx context.Context, scheID, resID string) error
}

type scheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepo{
		db: db,
	}
}

func (r *scheduleRepo) Create(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *scheduleRepo) Update(ctx context.Context, resID, scheID string, updateData map[string]interface{}) (*entity.Schedule, error) {
	var schedule entity.Schedule

	if err := r.db.WithContext(ctx).Model(&schedule).Clauses(clause.Returning{}).Where("id = ? AND restaurant_id = ?", scheID, resID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *scheduleRepo) Delete(ctx context.Context, scheID, resID string) error {
	return r.db.WithContext(ctx).Model(&entity.Schedule{}).Where("id = ? AND restaurant_id = ?", scheID, resID).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}).Error
}

func (r *scheduleRepo) FindByID(ctx context.Context, scheID, resID string) (*entity.Schedule, error) {
	var schedule entity.Schedule

	if err := r.db.
		WithContext(ctx).
		Preload("Restaurant", "is_deleted = ?", false).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", scheID, resID, false).
		First(&schedule).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *scheduleRepo) FindAllByResID(ctx context.Context, resID string) ([]*entity.Schedule, error) {
	var schedules []*entity.Schedule

	if err := r.db.
		WithContext(ctx).
		Model(&entity.Schedule{}).
		Preload("Restaurant", "is_deleted = ?", false).
		Where("restaurant_id = ? AND is_deleted = ?", resID, false).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

