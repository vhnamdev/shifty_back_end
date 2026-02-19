package repository

import (
	"context"
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	IsEmailExist(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, id, newPassword string) error
	GetRestaurantMembers(ctx context.Context, page int, limit int, restaurantID string, filter *dto.UserFilter) ([]*entity.User, int64, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Check email is exist or not
func (r *userRepo) IsEmailExist(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Function create user or we can call create account by email and password
func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Get user by email
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.
		WithContext(ctx).
		Preload("UserRestaurants").
		Preload("UserRestaurants.Restaurant").
		Preload("UserRestaurants.Position").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by Id
func (r *userRepo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).
		Preload("UserRestaurants").
		Preload("UserRestaurants.Restaurant").
		Preload("UserRestaurants.Position").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update user function
func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	result := r.db.WithContext(ctx).Model(user).Clauses(clause.Returning{}).Updates(user)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id, newPassword string) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("password", newPassword).Error
}

// Update user function
func (r *userRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"status":     false,
	}).Error
}

// Get all user
func (r *userRepo) GetRestaurantMembers(ctx context.Context, page, limit int, restaurantID string, filter *dto.UserFilter) ([]*entity.User, int64, error) {

	var users []*entity.User
	var total int64

	// Caculate offset
	offset := (page - 1) * limit

	// Create a query to search and limit it to only the user's restaurant.
	query := r.db.WithContext(ctx).Model(&entity.User{})

	query = query.
		Joins("JOIN user_restaurants ON user_restaurants.user_id = users.id").
		Where("user_restaurants.restaurant_id = ?", restaurantID)

	// Check if filter search is not nul
	if filter.Search != nil && *filter.Search != "" {

		// Prepare wildcard pattern for partial match
		searchPattern := "%" + *filter.Search + "%"

		// Apply case-insensitive search on both FullName and Email fields
		query = query.Where("(full_name ILIKE ? OR email ILIKE ?)", searchPattern, searchPattern)
	}

	// Check if filter search is not nul
	if filter.Position != nil && *filter.Position != "" {

		// Apply case-insensitive role field
		query = query.Where("user_restaurants.position_id = ?", *filter.Position)
	}

	// Count total records for pagination metadata
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Execute query with pagination, sorting, and eager loading
	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Preload("UserRestaurants", "restaurant_id = ?", restaurantID).
		Preload("UserRestaurants.Position").
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
