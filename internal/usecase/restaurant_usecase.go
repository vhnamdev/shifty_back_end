package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type RestaurantUseCase interface {
	Create(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	Update(ctx context.Context, userID string, restaurant *entity.Restaurant) error
	Delete(ctx context.Context, userID, resID string) error
	GetByID(ctx context.Context, userID, resID string) (*entity.Restaurant, error)
	GetMyRestaurants(ctx context.Context, userID string) ([]*entity.Restaurant, error)
}

type restaurantUseCase struct {
	RestaurantRepo     repository.RestaurantRepository
	UserRestaurantRepo repository.UserRestaurantRepository
}

func NewRestaurantUseCase(RestaurantRepo repository.RestaurantRepository, UserRestaurantRepo repository.UserRestaurantRepository) RestaurantUseCase {
	return &restaurantUseCase{
		RestaurantRepo:     RestaurantRepo,
		UserRestaurantRepo: UserRestaurantRepo,
	}
}

// Create restaurant
func (u *restaurantUseCase) Create(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {
	newRestaurant, err := u.RestaurantRepo.Create(ctx, restaurant)

	if err != nil {
		return nil, xerror.Internal("")
	}

	return newRestaurant, nil
}

func (u *restaurantUseCase) Update(ctx context.Context, userID string, restaurant *entity.Restaurant) (*entity.Restaurant, error) {

	isAuthority, err := u.UserRestaurantRepo.CheckAuthorityToUpdate(ctx, userID, restaurant.ID.String())

	if err != nil {
		return nil, xerror.Internal("Authority cannot be verified")
	}

	if !isAuthority {
		return nil, xerror.BadRequest("You are not allowed to update restaurant")
	}

	if err := u.RestaurantRepo.Update(ctx, restaurant); err != nil {
		return nil, xerror.Internal("Can not update restaurant")
	}

	return restaurant, nil
}

func (u *restaurantUseCase) Delete(ctx context.Context, userID, resID string) error {
	isAuthority, err := u.UserRestaurantRepo.CheckAuthorityToDelete(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Authority cannot be verified")
	}

	if !isAuthority {
		return xerror.BadRequest("You are not allowed to delete restaurant")
	}

	if err = u.RestaurantRepo.Delete(ctx, resID); err != nil {
		return xerror.Internal("Can not delete this restaurant")
	}

	return nil

}

func (u *restaurantUseCase) GetByID(ctx context.Context, userID, resID string) (*entity.Restaurant, error) {
	isAuthority, err := u.UserRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.BadRequest("Authority cannot be verified")
	}

	if !isAuthority {
		return nil, xerror.BadRequest("You can not allowed to see information of this restaurant")
	}

	restaurant, err := u.RestaurantRepo.GetByID(ctx, resID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Restaurant not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return restaurant, nil
}

