package usecase

import (
	"context"
	"log"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/mailer"
	"shifty-backend/pkg/uploader"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

type RestaurantUseCase interface {
	Create(ctx context.Context, userID string, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	Update(ctx context.Context, userID string, resID string, updateData map[string]interface{}) (*entity.Restaurant, error)
	UpdateImage(ctx context.Context, userID, resID, imageURL string) (*entity.Restaurant, error)
	Delete(ctx context.Context, userID, resID string) error
	GetByID(ctx context.Context, userID, resID string) (*entity.Restaurant, error)
	GetMyRestaurants(ctx context.Context, userID string) ([]*entity.Restaurant, error)
	CreateInviteCode(ctx context.Context, userID, email, resID, positionID string) error
	JoinRestaurant(ctx context.Context, userID, inviteCode string) error
}

type restaurantUseCase struct {
	transactor         repository.Transactor
	restaurantRepo     repository.RestaurantRepository
	userRestaurantRepo repository.UserRestaurantRepository
	positionRepo       repository.PositionRepository
	redisRepo          repository.RedisRepository
	userRepo           repository.UserRepository
	mailerService      mailer.EmailSender
	uploader           uploader.ImageUploader
}

func NewRestaurantUseCase(
	transactor repository.Transactor,
	restaurantRepo repository.RestaurantRepository,
	userRestaurantRepo repository.UserRestaurantRepository,
	positionRepo repository.PositionRepository,
	redisRepo repository.RedisRepository,
	userRepo repository.UserRepository,
	mailerService mailer.EmailSender,
	uploader uploader.ImageUploader) RestaurantUseCase {
	return &restaurantUseCase{
		transactor:         transactor,
		restaurantRepo:     restaurantRepo,
		userRestaurantRepo: userRestaurantRepo,
		positionRepo:       positionRepo,
		redisRepo:          redisRepo,
		userRepo:           userRepo,
		uploader:           uploader,
	}
}

// Create restaurant
func (u *restaurantUseCase) Create(ctx context.Context, userID string, restaurant *entity.Restaurant) (*entity.Restaurant, error) {

	// Parsed UserID
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return nil, xerror.BadRequest("Invalid user ID")
	}

	// Use transaction to handle 3 tasks in one time, If one of the three tasks reports an error, the steps will be undone
	err = u.transactor.WithTransaction(ctx, func(txCtx context.Context) error {

		// Task 1: create restaurant
		createdRes, err := u.restaurantRepo.Create(txCtx, restaurant)

		if err != nil {
			return err
		}

		// Task 2: create position owner
		ownerPosition := &entity.Position{
			Name:                constants.RoleOwner,
			Description:         constants.DescOwner,
			Rank:                constants.RankOwner,
			CanUpdateRestaurant: true,
			CanDeleteRestaurant: true,
			RestaurantID:        createdRes.ID,
		}

		createdPos, err := u.positionRepo.Create(txCtx, ownerPosition)
		if err != nil {
			return err
		}

		// Task 3: Create user restaurant
		userRes := &entity.UserRestaurant{
			UserID:       parsedUserID,
			RestaurantID: createdRes.ID,
			PositionID:   createdPos.ID,
		}

		_, err = u.userRestaurantRepo.Create(txCtx, userRes)

		if err != nil {
			return err
		}

		// If all task are completed, return nil
		return nil
	})

	// If one of the three tasks reports an error, the steps will be undone
	if err != nil {
		return nil, xerror.Internal("Failed to create restaurant and setup owner")
	}

	return restaurant, nil
}

// Update restaurant
func (u *restaurantUseCase) Update(ctx context.Context, userID string, resID string, updateData map[string]interface{}) (*entity.Restaurant, error) {

	// Check authority, if not owner return err
	isAuthority, err := u.userRestaurantRepo.CheckAuthorityToUpdate(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Authority cannot be verified")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to update restaurant")
	}

	// Check if user does not update any data, return restaurant's information
	if len(updateData) == 0 {
		return u.restaurantRepo.GetByID(ctx, resID)
	}

	// Update restaurant
	updatedRestaurant, err := u.restaurantRepo.Update(ctx, resID, updateData)

	if err != nil {
		return nil, xerror.Internal("Can not update restaurant")
	}

	return updatedRestaurant, nil
}

// Update restaurant' avatar
func (u *restaurantUseCase) UpdateImage(ctx context.Context, userID, resID, imageURL string) (*entity.Restaurant, error) {
	oldRestaurant, err := u.restaurantRepo.GetByID(ctx, resID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Restaurant can not found")
		}

		return nil, xerror.Internal("Database failed")
	}
	isAuthority, err := u.userRestaurantRepo.CheckAuthorityToUpdate(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Authority cannot be verified")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to update restaurant")
	}
	updatedRestaurant, err := u.restaurantRepo.UpdateImage(ctx, resID, imageURL)

	if err != nil {
		return nil, xerror.Internal("Can not update image")
	}

	if oldRestaurant.Avatar != "" && oldRestaurant.Avatar != imageURL {
		publicID := u.uploader.GetPublicIDFromURL(oldRestaurant.Avatar)
		if publicID != "" {
			return nil, xerror.BadRequest("Can not get publicID")
		}
		go func(pID string) {
			errDelete := u.uploader.DeleteImage(context.Background(), pID)

			if errDelete != nil {
				log.Printf("[CLEANUP ERROR]: Failed to delete old image from Cloudinary. PublicID: %s, Error: %v", pID, errDelete)
			} else {
				log.Printf("[CLEANUP SUCCESS]: Old image deleted. PublicID: %s", pID)
			}
		}(publicID)
	}
	return updatedRestaurant, nil
}

// Delete restaurant
func (u *restaurantUseCase) Delete(ctx context.Context, userID, resID string) error {

	// Similar to the update function, this function has to check authority of user
	isAuthority, err := u.userRestaurantRepo.CheckAuthorityToDelete(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Authority cannot be verified")
	}

	if !isAuthority {
		return xerror.Forbidden("You are not allowed to update restaurant")
	}

	// Delete
	if err = u.restaurantRepo.Delete(ctx, resID); err != nil {
		return xerror.Internal("Can not delete this restaurant")
	}

	return nil

}

// Get restaurant by userID and restaurant id, why need to check 2 IDs.
//
//	Because each user can work in many restaurants.
func (u *restaurantUseCase) GetByID(ctx context.Context, userID, resID string) (*entity.Restaurant, error) {

	// Check user is work in this restaurant or not
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.BadRequest("Authority cannot be verified")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to update restaurant")
	}

	// Get information
	restaurant, err := u.restaurantRepo.GetByID(ctx, resID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Restaurant not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return restaurant, nil
}

// Get all restaurants of this user
func (u *restaurantUseCase) GetMyRestaurants(ctx context.Context, userID string) ([]*entity.Restaurant, error) {
	restaurants, err := u.restaurantRepo.GetMyRestaurants(ctx, userID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Restaurants not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return restaurants, nil
}

func (u *restaurantUseCase) CreateInviteCode(ctx context.Context, userID, email, resID, positionID string) error {

	isAuthority, err := u.userRestaurantRepo.CheckAuthorityToInvite(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return xerror.Forbidden("You are not allowed to update restaurant")
	}

	inviteCode, err := utils.GenerateInviteCode(6)

	if err != nil {
		return xerror.Internal("Can not generate invite code")
	}

	err = u.redisRepo.SaveInviteCode(ctx, inviteCode, email, positionID, resID)

	if err != nil {
		return xerror.Internal("Save invite code failed")
	}

	position, err := u.positionRepo.FindByID(ctx, positionID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("Position can not found")
		}

		return xerror.Internal("Database failed")
	}

	restaurant, err := u.restaurantRepo.GetByID(ctx, resID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("Restaurant can not found")
		}

		return xerror.Internal("Database failed")
	}

	inviter, err := u.userRepo.GetByID(ctx, userID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("User can not found")
		}

		return xerror.Internal("Database failed")
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC SEND MAIL]: %v", r)
			}
		}()
		errMail := u.mailerService.SendInviteCode(email, inviter.FullName, restaurant.Name, position.Name, inviteCode)

		if errMail != nil {
			log.Printf("[MAIL ERROR] Target: %s | Restaurant: %s | Error: %v", email, restaurant.Name, errMail)
		} else {
			log.Printf("[MAIL SUCCESS] Invitation sent to: %s", email)
		}
	}()
	return nil
}

func (u *restaurantUseCase) JoinRestaurant(ctx context.Context, userID, inviteCode string) error {
	user, err := u.userRepo.GetByID(ctx, userID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("User not found")
		}
		return xerror.Internal("Database failed")
	}

	inviteData, err := u.redisRepo.VerifyInviteCode(ctx, user.Email, inviteCode)

	if err != nil {
		return err
	}
	resUUID, errRes := uuid.Parse(inviteData.RestaurantID)
	posUUID, errPos := uuid.Parse(inviteData.PositionID)
	userUUID, errUser := uuid.Parse(userID)
	if errRes != nil || errPos != nil || errUser != nil {
		return xerror.Internal("Invalid ID data")
	}
	exists, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userUUID.String(), resUUID.String())
	if err != nil {
		return xerror.Internal("Employee information verification error")
	}
	if exists {
		return xerror.BadRequest("You are already a member of this restaurant!")
	}
	userRes := &entity.UserRestaurant{
		RestaurantID: resUUID,
		PositionID:   posUUID,
		UserID:       userUUID,
		IsBanned:     false,
	}

	_, err = u.userRestaurantRepo.Create(ctx, userRes)

	if err != nil {
		return xerror.Internal("Can not join to restaurant")
	}

	return nil
}
