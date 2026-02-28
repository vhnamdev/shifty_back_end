package repository

import (
	"context"
	"shifty-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	Update(ctx context.Context, resID string, updateData map[string]interface{}) (*entity.Restaurant, error)
	UpdateImage(ctx context.Context, resID, imageURL string) (*entity.Restaurant, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Restaurant, error)
	GetMyRestaurants(ctx context.Context, userID string) ([]*entity.Restaurant, error)
}

type RestaurantRepo struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &RestaurantRepo{
		db: db,
	}
}

// Create restaurant
func (r *RestaurantRepo) Create(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {

	// Create and return the result
	result := r.db.WithContext(ctx).Clauses(clause.Returning{}).Create(restaurant)

	if result.Error != nil {
		return nil, result.Error
	}
	return restaurant, nil
}

// Update Restaurant
func (r *RestaurantRepo) Update(ctx context.Context, resID string, updateData map[string]interface{}) (*entity.Restaurant, error) {
	var updatedRestaurant entity.Restaurant

	// Update data and return new data
	if err := r.db.WithContext(ctx).
		Model(&updatedRestaurant).
		Clauses(clause.Returning{}).
		Where("id = ?", resID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &updatedRestaurant, nil
}

// Update image of restaurant
func (r *RestaurantRepo) UpdateImage(ctx context.Context, resID, imageURL string) (*entity.Restaurant, error) {
	var updatedRestaurant entity.Restaurant

	if err := r.db.WithContext(ctx).
		Model(&updatedRestaurant).
		Clauses(clause.Returning{}).
		Where("id = ?", resID).
		Update("avatar", imageURL).
		Error; err != nil {
		return nil, err
	}

	return &updatedRestaurant, nil

}

// Delete Restaurant, set IsDeleted equal true
func (r *RestaurantRepo) Delete(ctx context.Context, id string) error {
	db := Extract(ctx, r.db)
	// Set is_deleted equal true and status equal false
	return db.WithContext(ctx).Model(&entity.Restaurant{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"status":     false,
	}).Error
}

// Get Restaurant by ID
func (r *RestaurantRepo) GetByID(ctx context.Context, id string) (*entity.Restaurant, error) {
	var restaurant entity.Restaurant
	if err := r.db.
		WithContext(ctx).
		Preload("Positions", "is_deleted = ?", false).
		Preload("Users", "is_deleted = ? AND status = ?", false, true).
		Preload("Laws").
		Where("id = ?", id).
		First(&restaurant).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil

}

// Get User's Restaurants
func (r *RestaurantRepo) GetMyRestaurants(ctx context.Context, userID string) ([]*entity.Restaurant, error) {
	var restaurants []*entity.Restaurant
	if err := r.db.
		WithContext(ctx).
		Model(&entity.Restaurant{}).
		Joins("JOIN user_restaurants ON user_restaurants.restaurant_id = restaurants.id").
		Where("user_restaurants.user_id = ?", userID).
		Preload("Positions", "is_deleted = ?", false).
		Find(&restaurants).
		Error; err != nil {
		return nil, err
	}

	return restaurants, nil
}
