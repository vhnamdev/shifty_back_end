package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type ShiftAssignmentUseCase interface {
	Create(ctx context.Context, userID, resID, shiftID string, shiftAssignment *entity.ShiftAssignment) (*entity.ShiftAssignment, error)
	Update(ctx context.Context, userID, resID, shiftID, shiftAssignmentID string, updateData map[string]interface{}) (*entity.ShiftAssignment, error)
	Delete(ctx context.Context, userID, resID, shiftID, shiftAssignmentID string) error
	FindByID(ctx context.Context, userID, resID, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error)
	FindAllByShiftID(ctx context.Context, userID, resID, shiftID string) ([]*entity.ShiftAssignment, error)
}

type shiftAssignmentUseCase struct {
	shiftAssignmentRepo repository.ShiftAssignmentRepository
	userRestaurantRepo  repository.UserRestaurantRepository
}

func NewShiftAssignmentUseCase(shiftAssignmentRepo repository.ShiftAssignmentRepository, userRestaurantRepo repository.UserRestaurantRepository) ShiftAssignmentUseCase {
	return &shiftAssignmentUseCase{
		shiftAssignmentRepo: shiftAssignmentRepo,
		userRestaurantRepo:  userRestaurantRepo,
	}
}

func (u *shiftAssignmentUseCase) Create(ctx context.Context, userID, resID, shiftID string, shiftAssignment *entity.ShiftAssignment) (*entity.ShiftAssignment, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift assignment")
	}

	isExist, err := u.shiftAssignmentRepo.IsExistByUserAndShift(ctx, shiftID, shiftAssignment.UserID.String())

	if err != nil {
		return nil, xerror.Internal("Can not check shift assignment")
	}

	if isExist {
		return nil, xerror.BadRequest("Staff already has assignment in this shift")
	}

	newShiftAssignment, err := u.shiftAssignmentRepo.Create(ctx, shiftAssignment)

	if err != nil {
		return nil, xerror.Internal("Can not create shift assignment")
	}

	return newShiftAssignment, nil
}

func (u *shiftAssignmentUseCase) Update(ctx context.Context, userID, resID, shiftID, shiftAssignmentID string, updateData map[string]interface{}) (*entity.ShiftAssignment, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to update shift assignment")
	}

	updatedShiftAssignment, err := u.shiftAssignmentRepo.Update(ctx, updateData, shiftID, shiftAssignmentID)

	if err != nil {
		return nil, xerror.Internal("Can not update shift assignment")
	}

	return updatedShiftAssignment, nil
}

func (u *shiftAssignmentUseCase) Delete(ctx context.Context, userID, resID, shiftID, shiftAssignmentID string) error {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return xerror.Forbidden("You are not allowed to delete shift assignment")
	}

	if err := u.shiftAssignmentRepo.Delete(ctx, shiftID, shiftAssignmentID); err != nil {
		return xerror.Internal("Can not delete shift assignment")
	}
	return nil
}

func (u *shiftAssignmentUseCase) FindByID(ctx context.Context, userID, resID, shiftID, shiftAssignmentID string) (*entity.ShiftAssignment, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view shift assignment")
	}

	shiftAssignment, err := u.shiftAssignmentRepo.GetByID(ctx, shiftID, shiftAssignmentID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift assignment is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftAssignment, nil
}

func (u *shiftAssignmentUseCase) FindAllByShiftID(ctx context.Context, userID, resID, shiftID string) ([]*entity.ShiftAssignment, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view shift assignments")
	}

	shiftAssignments, err := u.shiftAssignmentRepo.GetAllByShiftID(ctx, shiftID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift assignments are not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftAssignments, nil
}
