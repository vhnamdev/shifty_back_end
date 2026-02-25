package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/xerror"
)

type UserRestaurantUseCase interface {
	UpdateStaffByManager(ctx context.Context, requestID, targetUserID, resID string, updateData map[string]interface{}) (*entity.UserRestaurant, error)
}
type userRestaurantUseCase struct {
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewUserRestaurantUseCase(userRestaurantRepo repository.UserRestaurantRepository) UserRestaurantUseCase {
	return &userRestaurantUseCase{
		userRestaurantRepo: userRestaurantRepo,
	}
}

// Update staff by manager
func (u *userRestaurantUseCase) UpdateStaffByManager(ctx context.Context, requestID, targetUserID, resID string, updateData map[string]interface{}) (*entity.UserRestaurant, error) {
	isAuthority, err := u.userRestaurantRepo.CheckAuthority(ctx, targetUserID, requestID, resID)

	if err != nil {
		return nil, err
	}
	if !isAuthority {
		return nil, xerror.BadRequest("You can not update this account")
	}
	updatedStaff, err := u.userRestaurantRepo.Update(ctx, targetUserID, resID, updateData)
	if err != nil {
		return nil, xerror.Internal("Can not update staff")
	}
	return updatedStaff, nil
}

