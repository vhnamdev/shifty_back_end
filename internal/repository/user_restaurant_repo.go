package repository

import (
	"context"
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRestaurantRepository interface {
	Create(ctx context.Context, userRes *entity.UserRestaurant) (*entity.UserRestaurant, error)
	CheckUserInRestaurant(ctx context.Context, userID, resID string) (bool, error)
	CheckAuthority(ctx context.Context, targetID, requestID, resID string) (bool, error)
	CheckAuthorityToUpdate(ctx context.Context, userID, resID string) (bool, error)
	CheckAuthorityToDelete(ctx context.Context, userID, resID string) (bool, error)
	CheckAuthorityToInvite(ctx context.Context, userID, resID string) (bool, error)
	Update(ctx context.Context, userID, resID string, updateData map[string]interface{}) (*entity.UserRestaurant, error)
}

type userRestaurantRepo struct {
	db *gorm.DB
}

func NewUserRestaurantRepository(db *gorm.DB) UserRestaurantRepository {
	return &userRestaurantRepo{
		db: db,
	}
}

// Create user restaurant
func (r *userRestaurantRepo) Create(ctx context.Context, userRes *entity.UserRestaurant) (*entity.UserRestaurant, error) {
	db := Extract(ctx, r.db)

	result := db.WithContext(ctx).Clauses(clause.Returning{}).Create(userRes)

	if result.Error != nil {
		return nil, result.Error
	}

	return userRes, nil
}

// Check user if user want to get informations the staff or members in their restaurant
func (r *userRestaurantRepo) CheckUserInRestaurant(ctx context.Context, userID, resID string) (bool, error) {

	var count int64
	if err := r.db.
		WithContext(ctx).
		Model(&entity.UserRestaurant{}).
		Where("user_id = ? AND restaurant_id = ?", userID, resID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}

// Check authorioty when the owner or managers want to update their members in their restaurant
func (r *userRestaurantRepo) CheckAuthority(ctx context.Context, targetID, requestID, resID string) (bool, error) {

	// Create an array to save requestor and the member who is updated
	var ranks []dto.UserRank

	// Select 2 collumn are user_id and rank in position, join positions table into user restaurants table to get the rank
	if err := r.db.
		WithContext(ctx).
		Select("user_restaurants.user_id, positions.rank").
		Joins("JOIN positions ON positions.id = user_restaurants.position_id").
		Where("user_restaurants.restaurant_id = ? AND user_restaurants.user_id IN ?", resID, []string{requestID, targetID}).
		Scan(&ranks).Error; err != nil {
		return false, err
	}

	if len(ranks) < 2 {
		return false, nil
	}

	var requestRank, targetRank int

	// Use for loop to take the information of rank of 2 members
	for _, item := range ranks {
		switch item.UserID {
		case requestID:
			requestRank = item.Rank
		case targetID:
			targetRank = item.Rank
		}
	}

	// Check if the requestor's rank is lower than member id, can update information. In this project
	// In this project, I've assigned the highest priority to the lowest priority, with number 1 being the highest.
	if requestRank < targetRank {
		return true, nil
	}
	return false, nil
}

func (r *userRestaurantRepo) CheckAuthorityToUpdate(ctx context.Context, userID, resID string) (bool, error) {
	var canUpdate bool

	if err := r.db.
		WithContext(ctx).
		Model(&entity.UserRestaurant{}).
		Select("positions.can_update_restaurant").
		Joins("JOIN positions ON user_restaurants.position_id = positions.id").
		Where("user_restaurants.user_id = ? AND user_restaurants.restaurant_id = ?", userID, resID).
		Scan(&canUpdate).Error; err != nil {
		return false, err
	}
	return canUpdate, nil

}

func (r *userRestaurantRepo) CheckAuthorityToDelete(ctx context.Context, userID, resID string) (bool, error) {
	var canDelete bool

	if err := r.db.
		WithContext(ctx).
		Model(&entity.UserRestaurant{}).
		Select("positions.can_delete_restaurant").
		Joins("JOIN positions ON user_restaurants.position_id = positions.id").
		Where("user_restaurants.user_id = ? AND user_restaurants.restaurant_id = ?", userID, resID).
		Scan(&canDelete).Error; err != nil {
		return false, err
	}
	return canDelete, nil

}

func (r *userRestaurantRepo) CheckAuthorityToInvite(ctx context.Context, userID, resID string) (bool, error) {
	var rank int

	if err := r.db.WithContext(ctx).
		Model(&entity.UserRestaurant{}).
		Select("positions.rank").
		Joins("JOIN positions ON user_restaurants.position_id = positions.id").
		Where("user_restaurants.user_id = ? AND user_restaurants.restaurant_id = ?", userID, resID).
		Scan(&rank).Error; err != nil {
		return false, err
	}
	if rank == 0 {
		return false, nil
	}

	return rank <= 3, nil

}

// Update position or ban member
func (r *userRestaurantRepo) Update(ctx context.Context, userID, resID string, updateData map[string]interface{}) (*entity.UserRestaurant, error) {

	//I have to use map to can update boolean field in database
	if err := r.db.WithContext(ctx).
		Model(&entity.UserRestaurant{}).
		Where("user_id = ? AND restaurant_id = ?", userID, resID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}

	var updatedRecord entity.UserRestaurant

	// Get the information again and return to usecase
	if err := r.db.WithContext(ctx).
		Preload("Position").
		Preload("Restaurant").
		Where("user_id = ? AND restaurant_id = ?", userID, resID).
		First(&updatedRecord).Error; err != nil {
		return nil, err
	}
	return &updatedRecord, nil
}
