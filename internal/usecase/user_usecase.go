package usecase

import (
	"context"
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type UserUseCase interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByID(ctx context.Context, ID string) (*entity.User, error)
	DeleteUser(ctx context.Context, ID string) error
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetRestaurantMembers(ctx context.Context, page, limit int, restaurantID string, filter *dto.UserFilter) ([]entity.User, int64, error)
	ValidateRestaurantAccess(ctx context.Context, currentUserID, TargetRestaurantID string) (bool, error)
}
type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

// Find user by email
func (u *userUseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	// Get user
	user, err := u.userRepo.GetByEmail(ctx, email)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("User not found")
		}
		return nil, xerror.Internal("Database error")
	}

	return user, nil
}

// Find user by ID
func (u *userUseCase) FindUserByID(ctx context.Context, ID string) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, ID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("User not found")
		}
		return nil, xerror.Internal("Database error")
	}

	return user, nil
}

// Delete User by ID, change IsDeleted field to true
func (u *userUseCase) DeleteUser(ctx context.Context, ID string) error {
	user, err := u.userRepo.GetByID(ctx, ID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("User not found")
		}
		return xerror.Internal("Database failed")
	}
	user.IsDeleted = true

	if err := u.userRepo.Delete(ctx, user); err != nil {
		return xerror.Internal("Failed to delete user")
	}
	return nil
}

// Update user
func (u *userUseCase) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, xerror.Internal("Can not update user")
	}

	updatedUser, err := u.userRepo.GetByID(ctx, user.ID.String())

	if err != nil {
		return nil, xerror.Internal("Update success but failed to retrieve new data")
	}
	return updatedUser, nil
}

// Get restaurant's members
func (u *userUseCase) GetRestaurantMembers(ctx context.Context, page, limit int, restaurantID string, filter *dto.UserFilter) ([]entity.User, int64, error) {
	users, total, err := u.userRepo.GetRestaurantMembers(ctx, page, limit, restaurantID, filter)

	if err != nil {
		return nil, 0, xerror.Internal("Can not get members")
	}

	return users, total, nil
}

// Validate access to restaurant
func (u *userUseCase) ValidateRestaurantAccess(ctx context.Context, currentUserID, targetRestaurantID string) (bool, error) {

	// Get user
	user, err := u.userRepo.GetByID(ctx, currentUserID)

	// Check is valid or not
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return false, xerror.NotFound("User not found")
		}
		return false, xerror.Internal("Database failed")
	}
	// If this is super admin
	if user.Role == "admin" {
		// accept access to restaurant
		return true, nil
	}
	// Check if user's restaurantID is equal target restaurantID
	if user.RestaurantID.String() == targetRestaurantID {

		return true, nil
	}
	// if not member in restaurant return false
	return false, nil
}
