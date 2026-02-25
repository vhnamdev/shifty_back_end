package usecase

import (
	"context"
	"log"
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/uploader"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"strings"
)

type UserUseCase interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByID(ctx context.Context, id string) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateImage(ctx context.Context, id, imageURl string) (*entity.User, error)
	GetRestaurantMembers(ctx context.Context, page, limit int, restaurantID string, filter *dto.UserFilter) ([]*entity.User, int64, error)
	ValidateRestaurantAccess(ctx context.Context, currentUserID, TargetRestaurantID string) (bool, error)
}
type userUseCase struct {
	userRepo           repository.UserRepository
	userRestaurantRepo repository.UserRestaurantRepository
	uploader           uploader.ImageUploader
	transactor         repository.Transactor
	restaurantRepo     repository.RestaurantRepository
}

func NewUserUseCase(userRepo repository.UserRepository, userRestaurantRepo repository.UserRestaurantRepository, uploader uploader.ImageUploader, transactor repository.Transactor, restaurantRepo repository.RestaurantRepository) UserUseCase {
	return &userUseCase{
		userRepo:           userRepo,
		userRestaurantRepo: userRestaurantRepo,
		uploader:           uploader,
		transactor:         transactor,
		restaurantRepo:     restaurantRepo,
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
func (u *userUseCase) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("User not found")
		}
		return nil, xerror.Internal("Database error")
	}

	return user, nil
}

// Delete User by ID, change IsDeleted field to true
func (u *userUseCase) DeleteUser(ctx context.Context, userID string) error {
	_, err := u.userRepo.GetByID(ctx, userID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("User not found")
		}
		return xerror.Internal("Database failed")
	}

	err = u.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		err := u.userRestaurantRepo.DeleteAllByUserID(txCtx, userID)

		if err != nil {
			return xerror.Internal("Can not delete user restaurant")
		}

		err = u.userRepo.Delete(txCtx, userID)

		if err != nil {
			return xerror.Internal("Can not delete user")
		}

		return nil
	})

	if err != nil {
		return xerror.Internal("Failed to delete user")
	}
	return nil
}

// Update user
func (u *userUseCase) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, xerror.Internal("Can not update user")
	}
	return user, nil
}

// Update user's avatar
func (u *userUseCase) UpdateImage(ctx context.Context, id, imageURl string) (*entity.User, error) {

	// Get user
	oldUser, err := u.userRepo.GetByID(ctx, id)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("User not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	// Update new image
	updatedUser, err := u.userRepo.UpdateImage(ctx, id, imageURl)
	if err != nil {
		return nil, err
	}

	// Check old avatar is exist or not, check new image is equal old avatar and check is this image already uploaded to Cloudinary?
	if oldUser.Avatar != "" && oldUser.Avatar != imageURl && strings.Contains(oldUser.Avatar, "cloudinary.com") {

		// Get public ID
		publicID := u.uploader.GetPublicIDFromURL(oldUser.Avatar)

		if publicID != "" {
			return nil, xerror.BadRequest("Can not get publicID")
		}
		go func(pID string) {

			// Delete old image
			errDelete := u.uploader.DeleteImage(ctx, pID)

			if errDelete != nil {
				log.Printf("[CLEANUP ERROR]: Failed to delete old image from Cloudinary. PublicID: %s, Error: %v", pID, errDelete)
			} else {
				log.Printf("[CLEANUP SUCCESS]: Old image deleted. PublicID: %s", pID)
			}
		}(publicID)
	}
	return updatedUser, nil
}

// Get restaurant's members
func (u *userUseCase) GetRestaurantMembers(ctx context.Context, page, limit int, restaurantID string, filter *dto.UserFilter) ([]*entity.User, int64, error) {
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
	if user.Role == constants.Admin {
		// accept access to restaurant
		return true, nil
	}
	// Check if user's restaurantID is equal target restaurantID
	isMember, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, user.ID.String(), targetRestaurantID)

	if err != nil {
		return false, xerror.Internal("Check membership failed")
	}
	if isMember {
		return true, nil
	}
	// if not member in restaurant return false
	return false, nil
}
