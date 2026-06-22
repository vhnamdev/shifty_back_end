package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShiftRequestRepository interface {
	Create(ctx context.Context, shiftRequest *entity.ShiftRequest) (*entity.ShiftRequest, error)
	Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftRequestID string) (*entity.ShiftRequest, error)
	Delete(ctx context.Context, shiftID, shiftRequestID string) error
	GetByID(ctx context.Context, shiftID, shiftRequestID string) (*entity.ShiftRequest, error)
	GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftRequest, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*entity.ShiftRequest, error)
}

type shiftRequestRepo struct {
	db *gorm.DB
}

func NewShiftRequestRepository(db *gorm.DB) ShiftRequestRepository {
	return &shiftRequestRepo{
		db: db,
	}
}

func (r *shiftRequestRepo) Create(ctx context.Context, shiftRequest *entity.ShiftRequest) (*entity.ShiftRequest, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(shiftRequest).Error; err != nil {
		return nil, err
	}
	return shiftRequest, nil
}

func (r *shiftRequestRepo) Update(ctx context.Context, updateData map[string]interface{}, shiftID, shiftRequestID string) (*entity.ShiftRequest, error) {
	var shiftRequest entity.ShiftRequest

	if err := r.db.
		WithContext(ctx).
		Model(&shiftRequest).
		Clauses(clause.Returning{}).
		Where("id = ? AND shift_id = ? AND is_deleted = ?", shiftRequestID, shiftID, false).
		Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &shiftRequest, nil
}

func (r *shiftRequestRepo) Delete(ctx context.Context, shiftID, shiftRequestID string) error {
	return r.db.
		WithContext(ctx).
		Model(&entity.ShiftRequest{}).
		Where("id = ? AND shift_id = ? AND is_deleted = ?", shiftRequestID, shiftID, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *shiftRequestRepo) GetByID(ctx context.Context, shiftID, shiftRequestID string) (*entity.ShiftRequest, error) {
	var shiftRequest entity.ShiftRequest

	if err := r.db.
		WithContext(ctx).
		Preload("User", "is_deleted = ?", false).
		Preload("Shift", "is_deleted = ?", false).
		Preload("Position", "is_deleted = ?", false).
		Where("id = ? AND shift_id = ? AND is_deleted = ?", shiftRequestID, shiftID, false).
		First(&shiftRequest).Error; err != nil {
		return nil, err
	}
	return &shiftRequest, nil
}

func (r *shiftRequestRepo) GetAllByShiftID(ctx context.Context, shiftID string) ([]*entity.ShiftRequest, error) {
	var shiftRequests []*entity.ShiftRequest

	if err := r.db.
		WithContext(ctx).
		Model(&entity.ShiftRequest{}).
		Preload("User", "is_deleted = ?", false).
		Preload("Shift", "is_deleted = ?", false).
		Preload("Position", "is_deleted = ?", false).
		Where("shift_id = ? AND is_deleted = ?", shiftID, false).
		Find(&shiftRequests).Error; err != nil {
		return nil, err
	}
	return shiftRequests, nil
}

func (r *shiftRequestRepo) GetAllByUserID(ctx context.Context, userID string) ([]*entity.ShiftRequest, error) {
	var shiftRequests []*entity.ShiftRequest

	if err := r.db.
		WithContext(ctx).
		Model(&entity.ShiftRequest{}).
		Preload("User", "is_deleted = ?", false).
		Preload("Shift", "is_deleted = ?", false).
		Preload("Position", "is_deleted = ?", false).
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Find(&shiftRequests).Error; err != nil {
		return nil, err
	}
	return shiftRequests, nil
}
