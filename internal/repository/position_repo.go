package repository

import (
	"context"
	"shifty-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PositionRepository interface {
	Create(ctx context.Context, position *entity.Position) (*entity.Position, error)
	FindByID(ctx context.Context, ID string) (*entity.Position, error)
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

func (r *positionRepo) FindByID(ctx context.Context, ID string) (*entity.Position, error) {
	var position entity.Position

	if err := r.db.
		WithContext(ctx).
		Model(&entity.Position{}).
		Where("id = ?", ID).
		First(&position).
		Error; err != nil {
		return nil, err
	}

	return &position, nil
}
