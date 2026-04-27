package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShiftRepository interface {
	Create(ctx context.Context, shift *entity.Shift) (*entity.Shift, error)
	Update(ctx context.Context, shiftID, scheID string, updateData map[string]interface{}) (*entity.Shift, error)
	Delete(ctx context.Context, scheID, shiftID string) error
	FindByID(ctx context.Context, scheID, shiftID string) (*entity.Shift, error)
	FindAllByScheduleID(ctx context.Context, scheID string) ([]*entity.Shift, error)
}

type shiftRepo struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepo{
		db: db,
	}
}

func (r *shiftRepo) Create(ctx context.Context, shift *entity.Shift) (*entity.Shift, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(shift).Error; err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *shiftRepo) Update(ctx context.Context, shiftID, scheID string, updateData map[string]interface{}) (*entity.Shift, error) {
	var shift entity.Shift

	if err := r.db.
		WithContext(ctx).
		Model(&shift).
		Clauses(clause.Returning{}).
		Where("id = ? AND schedule_id = ? AND is_deleted = ?", shift, scheID, false).
		Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &shift, nil
}

func (r *shiftRepo) Delete(ctx context.Context, scheID, shiftID string) error {
	return r.db.WithContext(ctx).Model(&entity.Shift{}).Where("id = ? AND schedule_id = ?", shiftID, scheID).Updates(map[string]interface{}{
		"is_deleted": false,
		"deleted_at": time.Now(),
	}).Error
}

func (r *shiftRepo) FindByID(ctx context.Context, scheID, shiftID string) (*entity.Shift, error) {
	var shift entity.Shift

	if err := r.db.WithContext(ctx).Preload("Schedule", "is_deleted = ?", false).Where("id = ? AND schedule_id = ? AND is_deleted = ?", shift, scheID, false).First(&shift).Error; err != nil {
		return nil, err
	}

	return &shift, nil
}

func (r *shiftRepo) FindAllByScheduleID(ctx context.Context, scheID string) ([]*entity.Shift, error) {
	var shifts []*entity.Shift

	if err := r.db.WithContext(ctx).Model(&entity.Shift{}).Preload("Schedule", "is_deleted = ?", false).Where("schedule_id = ? AND is_deleted = ?", scheID, false).Find(&shifts).Error; err != nil {
		return nil, err
	}

	return shifts, nil
}
