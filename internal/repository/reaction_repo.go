package repository

import (
	"context"
	"shifty-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReactionRepository interface {
	Create(ctx context.Context, reaction *entity.Reaction) (*entity.Reaction, error)
	UpdateType(ctx context.Context, reactionID, reactionType string) (*entity.Reaction, error)
	Delete(ctx context.Context, reactionID string) error
	GetByPostAndAuthor(ctx context.Context, postID, authorID string) (*entity.Reaction, error)
	CountByPostID(ctx context.Context, postID string) (map[string]int, int, error)
}

type reactionRepo struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) ReactionRepository {
	return &reactionRepo{db: db}
}

func (r *reactionRepo) Create(ctx context.Context, reaction *entity.Reaction) (*entity.Reaction, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(reaction).Error; err != nil {
		return nil, err
	}
	return reaction, nil
}

func (r *reactionRepo) UpdateType(ctx context.Context, reactionID, reactionType string) (*entity.Reaction, error) {
	var reaction entity.Reaction

	if err := r.db.
		WithContext(ctx).
		Model(&reaction).
		Clauses(clause.Returning{}).
		Where("id = ?", reactionID).
		Updates(map[string]interface{}{"type": reactionType}).Error; err != nil {
		return nil, err
	}

	return &reaction, nil
}

func (r *reactionRepo) Delete(ctx context.Context, reactionID string) error {
	return r.db.WithContext(ctx).Where("id = ?", reactionID).Delete(&entity.Reaction{}).Error
}

func (r *reactionRepo) GetByPostAndAuthor(ctx context.Context, postID, authorID string) (*entity.Reaction, error) {
	var reaction entity.Reaction

	if err := r.db.
		WithContext(ctx).
		Where("post_id = ? AND author_id = ?", postID, authorID).
		First(&reaction).Error; err != nil {
		return nil, err
	}

	return &reaction, nil
}

func (r *reactionRepo) CountByPostID(ctx context.Context, postID string) (map[string]int, int, error) {
	type reactionCount struct {
		Type  string
		Count int
	}

	var rows []reactionCount
	if err := r.db.
		WithContext(ctx).
		Model(&entity.Reaction{}).
		Select("type, COUNT(*) as count").
		Where("post_id = ?", postID).
		Group("type").
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	counts := make(map[string]int, len(rows))
	total := 0
	for _, row := range rows {
		counts[row.Type] = row.Count
		total += row.Count
	}

	return counts, total, nil
}
