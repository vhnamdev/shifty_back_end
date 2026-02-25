package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type PositionUseCase interface {
	Create(ctx context.Context, position *entity.Position) (*entity.Position, error)
}

type positionUseCase struct {
	positionRepo       repository.PositionRepository
	userRestaurantRepo repository.UserRestaurantRepository
	transactor         repository.Transactor
}

func NewPositionUseCase(positionRepo repository.PositionRepository, userRestaurantRepo repository.UserRestaurantRepository, transactor repository.Transactor) PositionUseCase {
	return &positionUseCase{
		positionRepo:       positionRepo,
		userRestaurantRepo: userRestaurantRepo,
		transactor:         transactor,
	}
}

func (u *positionUseCase) Create(ctx context.Context, position *entity.Position) (*entity.Position, error) {
	result, err := u.positionRepo.Create(ctx, position)

	if err != nil {
		return nil, xerror.Internal("Can not create position")
	}

	return result, nil
}

func (u *positionUseCase) Update(ctx context.Context, posID, userID, resID string, updateData map[string]interface{}) (*entity.Position, error) {
	isAuthority, err := u.userRestaurantRepo.CheckAuthorityToUpdate(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You can not allowed to update position")
	}

	updatePosition, err := u.positionRepo.Update(ctx, posID, updateData)

	if err != nil {
		return nil, xerror.Internal("Can not update position")
	}

	return updatePosition, nil
}

func (u *positionUseCase) Delete(ctx context.Context, userID, resID, posID string) error {
	isAuthority, err := u.userRestaurantRepo.CheckAuthorityToDelete(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return xerror.Forbidden("You can not allowed to delete position")
	}

	err = u.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		err = u.userRestaurantRepo.SetPositionNull(txCtx, posID, resID)

		if err != nil {
			return xerror.Internal("Can not set position is user restaurant to null")
		}

		err = u.positionRepo.Delete(txCtx, posID, resID)

		if err != nil {
			return xerror.Internal("Can not delete this position")
		}

		return nil
	})

	if err != nil {
		return xerror.Internal("Failed to delete position")
	}

	return nil
}

func (u *positionUseCase) FindByID(ctx context.Context, posID, userID, resID string) (*entity.Position, error) {
	isMember, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Failed to verify membership")
	}

	if !isMember {
		return nil, xerror.Forbidden("You are not allowed of this restaurant")
	}

	position, err := u.positionRepo.FindByID(ctx, posID, resID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Position not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return position, nil
}

func (u *positionUseCase) GetAllByRestaurantID(ctx context.Context, resID, userID string) ([]*entity.Position, error) {
	isMember, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Failed to verify membership")
	}

	if !isMember {
		return nil, xerror.Forbidden("You are not allowed of this restaurant")
	}

	positions, err := u.positionRepo.GetAllByRestaurantID(ctx, resID)

	if err != nil {
		return nil, xerror.Internal("Can not get positions")
	}

	return positions, nil

}
