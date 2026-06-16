package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShiftRequirementRepository interface {
	Create(ctx context.Context, shiftRequirement *entity.ShiftRequirement) (*entity.ShiftRequirement, error)
	Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error)
	Delete(ctx context.Context, shiftID, shiftRequireID string) error
	GetByID(ctx context.Context, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error)
	GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftRequirement, error)
}

type shiftRequirementRepo struct {
	db *gorm.DB
}

func NewShiftRequirementRepository(db *gorm.DB) ShiftRequirementRepository {
	return &shiftRequirementRepo{
		db: db,
	}
}

func (r *shiftRequirementRepo) Create(ctx context.Context, shiftRequirement *entity.ShiftRequirement) (*entity.ShiftRequirement, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(shiftRequirement).Error; err != nil {
		return nil, err
	}
	return shiftRequirement, nil
}

func (r *shiftRequirementRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error) {
	var shiftRequire entity.ShiftRequirement

	if err := r.db.
		WithContext(ctx).
		Model(&shiftRequire).
		Clauses(clause.Returning{}).
		Where("id = ? AND shift_id = ?", shiftRequireID, shiftID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &shiftRequire, nil
}

func (r *shiftRequirementRepo) Delete(ctx context.Context, shiftID, shiftRequireID string) error {
	return r.db.WithContext(ctx).Model(&entity.ShiftRequirement{}).Where("id = ? AND shift_id = ?", shiftRequireID, shiftID).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}).Error
}

func (r *shiftRequirementRepo) GetByID(ctx context.Context, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error) {
	var shiftRequire entity.ShiftRequirement

	if err := r.db.
		WithContext(ctx).
		Preload("Shift", "is_deleted = ?", false).
		Where("id = ? AND shift_id = ? AND is_deleted = ?", shiftRequireID, shiftID, false).
		First(&shiftRequire).Error; err != nil {
		return nil, err
	}
	return &shiftRequire, nil
}

func (r *shiftRequirementRepo) GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftRequirement, error) {
	var shiftRequires []*entity.ShiftRequirement

	if err := r.db.
		WithContext(ctx).
		Model(&entity.ShiftRequirement{}).
		Preload("Shift", "is_deleted = ?", false).
		Where("shift_id = ? AND is_deleted = ?", shiftID, false).
		Find(&shiftRequires).Error; err != nil {
		return nil, err
	}
	return shiftRequires, nil
}
