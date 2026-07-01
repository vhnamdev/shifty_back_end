package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, updateData map[string]interface{}, commentID, postID string) (*entity.Comment, error)
	Delete(ctx context.Context, commentID, postID string) error
	GetByID(ctx context.Context, commentID, postID string) (*entity.Comment, error)
	GetAllByPostID(ctx context.Context, postID string, page, limit int) ([]*entity.Comment, int64, error)
}

type commentRepo struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepo{
		db: db,
	}
}

func (r *commentRepo) Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *commentRepo) Update(ctx context.Context, updateData map[string]interface{}, commentID, postID string) (*entity.Comment, error) {
	var comment entity.Comment

	if err := r.db.
		WithContext(ctx).
		Model(&comment).
		Clauses(clause.Returning{}).
		Where("id = ? AND post_id = ? AND is_deleted = ?", commentID, postID, false).
		Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepo) Delete(ctx context.Context, commentID, postID string) error {
	return r.db.
		WithContext(ctx).
		Model(&entity.Comment{}).
		Where("id = ? AND post_id = ? AND is_deleted = ?", commentID, postID, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *commentRepo) GetByID(ctx context.Context, commentID, postID string) (*entity.Comment, error) {
	var comment entity.Comment

	if err := r.db.
		WithContext(ctx).
		Preload("User", "is_deleted = ?", false).
		Preload("Replies", "is_deleted = ?", false).
		Where("id = ? AND post_id = ? AND is_deleted = ?", commentID, postID, false).
		First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepo) GetAllByPostID(ctx context.Context, postID string, page, limit int) ([]*entity.Comment, int64, error) {
	var comments []*entity.Comment
	var total int64

	offset := (page - 1) * limit

	query := r.db.
		WithContext(ctx).
		Model(&entity.Comment{}).
		Where("post_id = ? AND parent_id IS NULL AND is_deleted = ?", postID, false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("User", "is_deleted = ?", false).
		Preload("Replies", "is_deleted = ?", false).
		Limit(limit).
		Offset(offset).
		Order("created_at asc").
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
