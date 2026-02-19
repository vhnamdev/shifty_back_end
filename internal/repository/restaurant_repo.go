package repository

import (
	"context"
	"shifty-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	Update(ctx context.Context, restaurant *entity.Restaurant) error
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
func (r *RestaurantRepo) Update(ctx context.Context, restaurant *entity.Restaurant) error {

	// Update data and return new data
	result := r.db.WithContext(ctx).Clauses(clause.Returning{}).Updates(restaurant)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete Restaurant, set IsDeleted equal true
func (r *RestaurantRepo) Delete(ctx context.Context, id string) error {

	// Set is_deleted equal true and status equal false
	return r.db.WithContext(ctx).Model(&entity.Restaurant{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"status":     false,
	}).Error
}

// Get Restaurant by ID
func (r *RestaurantRepo) GetByID(ctx context.Context, id string) (*entity.Restaurant, error) {
	var restaurant entity.Restaurant
	if err := r.db.
		WithContext(ctx).
		Preload("Positions").
		Preload("Users").
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
		Joins("Join users ON users.restaurant_id = restaurants.id").
		Where("users.id = ?", userID).
		Preload("Positions").
		Preload("Laws").
		Preload("Users").
		Find(&restaurants).
		Error; err != nil {
		return nil, err
	}

	return restaurants, nil
}
