package repository

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) (*entity.Post, error)
	Update(ctx context.Context, updateData map[string]interface{}, postID, restaurantID string) (*entity.Post, error)
	Delete(ctx context.Context, postID, restaurantID string) error
	GetByID(ctx context.Context, postID, restaurantID string) (*entity.Post, error)
	GetAllByRestaurantID(ctx context.Context, restaurantID string, page, limit int) ([]*entity.Post, int64, error)
}

type postRepo struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepo{
		db: db,
	}
}

func (r *postRepo) Create(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepo) Update(ctx context.Context, updateData map[string]interface{}, postID, restaurantID string) (*entity.Post, error) {
	var post entity.Post

	if err := r.db.
		WithContext(ctx).
		Model(&post).
		Clauses(clause.Returning{}).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", postID, restaurantID, false).
		Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) Delete(ctx context.Context, postID, restaurantID string) error {
	return r.db.
		WithContext(ctx).
		Model(&entity.Post{}).
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", postID, restaurantID, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *postRepo) GetByID(ctx context.Context, postID, restaurantID string) (*entity.Post, error) {
	var post entity.Post

	if err := r.db.
		WithContext(ctx).
		Preload("Restaurant", "is_deleted = ?", false).
		Preload("User", "is_deleted = ?", false).
		Preload("Reactions").
		Preload("Comments").
		Where("id = ? AND restaurant_id = ? AND is_deleted = ?", postID, restaurantID, false).
		First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) GetAllByRestaurantID(ctx context.Context, restaurantID string, page, limit int) ([]*entity.Post, int64, error) {
	var posts []*entity.Post
	var total int64

	offset := (page - 1) * limit

	query := r.db.
		WithContext(ctx).
		Model(&entity.Post{}).
		Where("restaurant_id = ? AND is_deleted = ?", restaurantID, false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Restaurant", "is_deleted = ?", false).
		Preload("User", "is_deleted = ?", false).
		Preload("Reactions").
		Preload("Comments").
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
