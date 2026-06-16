package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type ShiftRequirementUseCase interface {
	Create(ctx context.Context, userID, resID, shiftID string, shiftRequire *entity.ShiftRequirement) (*entity.ShiftRequirement, error)
	Update(ctx context.Context, userID, resID, shiftID, shiftRequireID string, updateData map[string]interface{}) (*entity.ShiftRequirement, error)
	Delete(ctx context.Context, userID, resID, shiftID, shiftRequireID string) error
	FindByID(ctx context.Context, userID, resID, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error)
	FindAllByShiftID(ctx context.Context, userID, resID, shiftID string) ([]*entity.ShiftRequirement, error)
}

type shiftRequirementUseCase struct {
	shiftRequirementRepo repository.ShiftRequirementRepository
	userRestaurantRepo   repository.UserRestaurantRepository
}

func NewShiftRequirementUseCase(shiftRequirementRepo repository.ShiftRequirementRepository, userRestaurantRepo repository.UserRestaurantRepository) ShiftRequirementUseCase {
	return &shiftRequirementUseCase{
		shiftRequirementRepo: shiftRequirementRepo,
		userRestaurantRepo:   userRestaurantRepo,
	}
}

func (u *shiftRequirementUseCase) Create(ctx context.Context, userID, resID, shiftID string, shiftRequire *entity.ShiftRequirement) (*entity.ShiftRequirement, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift")
	}

	newShiftRequire, err := u.shiftRequirementRepo.Create(ctx, shiftRequire)

	if err != nil {
		return nil, xerror.Internal("Can not create shift")
	}

	return newShiftRequire, nil
}

func (u *shiftRequirementUseCase) Update(ctx context.Context, userID, resID, shiftID, shiftRequireID string, updateData map[string]interface{}) (*entity.ShiftRequirement, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift")
	}

	updatedShiftRequire, err := u.shiftRequirementRepo.Update(ctx, updateData, shiftID, shiftRequireID)

	if err != nil {
		return nil, xerror.Internal("Can not update shift")
	}

	return updatedShiftRequire, nil
}

func (u *shiftRequirementUseCase) Delete(ctx context.Context, userID, resID, shiftID, shiftRequireID string) error {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return xerror.Forbidden("You are not allowed to create shift")
	}

	return u.shiftRequirementRepo.Delete(ctx, shiftID, shiftRequireID)
}

func (u *shiftRequirementUseCase) FindByID(ctx context.Context, userID, resID, shiftID, shiftRequireID string) (*entity.ShiftRequirement, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift")
	}

	shiftRequire, err := u.shiftRequirementRepo.GetByID(ctx, shiftID, shiftRequireID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift requirement is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftRequire, nil
}

func (u *shiftRequirementUseCase) FindAllByShiftID(ctx context.Context, userID, resID, shiftID string) ([]*entity.ShiftRequirement, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift")
	}

	shiftRequires, err := u.shiftRequirementRepo.GetAllByShiftID(ctx, shiftID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift requirements are not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftRequires, nil
}
