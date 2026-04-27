package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type ShiftUseCase interface {
	Create(ctx context.Context, userID, resID string, shift *entity.Shift) (*entity.Shift, error)
	Update(ctx context.Context, userID, resID, scheID, shiftID string, updateData map[string]interface{}) (*entity.Shift, error)
	Delete(ctx context.Context, userID, resID, scheID, shiftID string) error
	FindByID(ctx context.Context, userID, resID, scheID, shiftID string) (*entity.Shift, error)
	FindAllByScheduleID(ctx context.Context, userID, resID, scheID string) ([]*entity.Shift, error)
}

type shiftUseCase struct {
	shiftRepo          repository.ShiftRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewShiftUseCase(shiftRepo repository.ShiftRepository, userRestaurantRepo repository.UserRestaurantRepository) ShiftUseCase {
	return &shiftUseCase{
		shiftRepo:          shiftRepo,
		userRestaurantRepo: userRestaurantRepo,
	}
}

func (u *shiftUseCase) Create(ctx context.Context, userID, resID string, shift *entity.Shift) (*entity.Shift, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift")
	}

	newShift, err := u.shiftRepo.Create(ctx, shift)

	if err != nil {
		return nil, xerror.Internal("Can not create shift")
	}

	return newShift, nil
}

func (u *shiftUseCase) Update(ctx context.Context, userID, resID, scheID, shiftID string, updateData map[string]interface{}) (*entity.Shift, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to update shift")
	}

	updatedShift, err := u.shiftRepo.Update(ctx, shiftID, scheID, updateData)

	if err != nil {
		return nil, xerror.Internal("Can not update shift")
	}

	return updatedShift, nil
}

func (u *shiftUseCase) Delete(ctx context.Context, userID, resID, scheID, shiftID string) error {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return xerror.Forbidden("You are not allowed to delete shift")
	}

	if err := u.shiftRepo.Delete(ctx, scheID, shiftID); err != nil {
		return xerror.Internal("Can not delete shift")
	}
	return nil
}

func (u *shiftUseCase) FindByID(ctx context.Context, userID, resID, scheID, shiftID string) (*entity.Shift, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to access shift in this restaurant")
	}

	shift, err := u.shiftRepo.FindByID(ctx, scheID, shiftID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shift, nil
}

func (u *shiftUseCase) FindAllByScheduleID(ctx context.Context, userID, resID, scheID string) ([]*entity.Shift, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to access shift in this restaurant")
	}

	shifts, err := u.shiftRepo.FindAllByScheduleID(ctx, scheID)

	if err != nil {
		return nil, err
	}

	return shifts, nil
}
