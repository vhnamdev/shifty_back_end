package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShiftAssignmentRepository interface {
	Create(ctx context.Context, shiftAssignment *entity.ShiftAssignment) (*entity.ShiftAssignment, error)
	Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error)
	Delete(ctx context.Context, shiftID, shiftAssignmentID string) error
	GetByID(ctx context.Context, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error)
	GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftAssignment, error)
	IsExistByUserAndShift(ctx context.Context, shiftID, userID string) (bool, error)
}

type shiftAssignmentRepo struct {
	db *gorm.DB
}

func NewShiftAssignmentRepository(db *gorm.DB) ShiftAssignmentRepository {
	return &shiftAssignmentRepo{
		db: db,
	}
}

func (r *shiftAssignmentRepo) Create(ctx context.Context, shiftAssignment *entity.ShiftAssignment) (*entity.ShiftAssignment, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(shiftAssignment).Error; err != nil {
		return nil, err
	}
	return shiftAssignment, nil
}

func (r *shiftAssignmentRepo) IsExistByUserAndShift(ctx context.Context, shiftID, userID string) (bool, error) {
	var count int64

	if err := r.db.
		WithContext(ctx).
		Model(&entity.ShiftAssignment{}).
		Where("shift_id = ? AND user_id = ? AND is_deleted = ?", shiftID, userID, false).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *shiftAssignmentRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error) {
	var shiftAssignment entity.ShiftAssignment

	if err := r.db.
		WithContext(ctx).
		Model(&shiftAssignment).
		Clauses(clause.Returning{}).
		Where("id = ? AND shift_id = ?", shiftAssignmentID, shiftID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &shiftAssignment, nil
}

func (r *shiftAssignmentRepo) Delete(ctx context.Context, shiftID, shiftAssignmentID string) error {
	if err := r.db.
		WithContext(ctx).
		Model(&entity.ShiftAssignment{}).
		Where("id = ? AND shift_id = ? AND is_deleted = ?", shiftAssignmentID, shiftID, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *shiftAssignmentRepo) GetByID(ctx context.Context, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error) {
	var shiftAssignment entity.ShiftAssignment

	if err := r.db.
		WithContext(ctx).
		Preload("User", "is_deleted = ?", false).
		Preload("Shift", "is_deleted = ?", false).
		Preload("Position", "is_deleted = ?", false).
		Where("id = ? AND shift_id = ? AND is_deleted = ?", shiftAssignmentID, shiftID, false).
		First(&shiftAssignment).Error; err != nil {
		return nil, err
	}
	return &shiftAssignment, nil
}

func (r *shiftAssignmentRepo) GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftAssignment, error) {
	var shiftAssignments []*entity.ShiftAssignment

	if err := r.db.
		WithContext(ctx).
		Model(&entity.ShiftAssignment{}).
		Preload("User", "is_deleted = ?", false).
		Preload("Shift", "is_deleted = ?", false).
		Preload("Position", "is_deleted = ?", false).
		Where("shift_id = ? AND is_deleted = ?", shiftID, false).
		Find(&shiftAssignments).Error; err != nil {
		return nil, err
	}
	return shiftAssignments, nil
}
