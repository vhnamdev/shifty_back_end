package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PositionRepository interface {
	Create(ctx context.Context, position *entity.Position) (*entity.Position, error)
	FindByID(ctx context.Context, posID, resID string) (*entity.Position, error)
	GetAllByRestaurantID(ctx context.Context, resID string) ([]*entity.Position, error)
	Delete(ctx context.Context, posID, resID string) error
	DeleteAllByRestaurantID(ctx context.Context, resID string) error
	Update(ctx context.Context, posID, resID string, updateData map[string]interface{}) (*entity.Position, error)
}

type positionRepo struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positionRepo{
		db: db,
	}
}

func (r *positionRepo) Create(ctx context.Context, position *entity.Position) (*entity.Position, error) {
	db := Extract(ctx, r.db)

	result := db.WithContext(ctx).Clauses(clause.Returning{}).Create(position)

	if result.Error != nil {
		return nil, result.Error
	}

	return position, nil
}

func (r *positionRepo) Update(ctx context.Context, posID, resID string, updateData map[string]interface{}) (*entity.Position, error) {
	var updatedPosition entity.Position

	result := r.db.WithContext(ctx).Model(&updatedPosition).Clauses(clause.Returning{}).Where("id = ? AND restaurant_id = ?", posID, resID).Updates(updateData)

	if result.Error != nil {
		return nil, result.Error
	}

	return &updatedPosition, nil
}

func (r *positionRepo) Delete(ctx context.Context, posID, resID string) error {

	db := Extract(ctx, r.db)
	return db.WithContext(ctx).Model(&entity.Position{}).Where("id = ? AND restaurant_id = ?", posID, resID).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}).Error
}

func (r *positionRepo) DeleteAllByRestaurantID(ctx context.Context, resID string) error {
	db := Extract(ctx, r.db)

	return db.WithContext(ctx).Model(&entity.Position{}).Where("restaurant_id = ?", resID).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}).Error
}

func (r *positionRepo) FindByID(ctx context.Context, posID, resID string) (*entity.Position, error) {
	var position entity.Position

	if err := r.db.
		WithContext(ctx).
		Model(&entity.Position{}).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", posID, resID, false).
		First(&position).
		Error; err != nil {
		return nil, err
	}

	return &position, nil
}

func (r *positionRepo) GetAllByRestaurantID(ctx context.Context, resID string) ([]*entity.Position, error) {
	var positions []*entity.Position

	if err := r.db.
		WithContext(ctx).
		Where("restaurant_id = ? AND is_deleted = ?", resID, false).
		Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}
