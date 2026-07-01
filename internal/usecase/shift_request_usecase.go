package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type ShiftRequestUseCase interface {
	Create(ctx context.Context, userID, resID, shiftID string, shiftRequest *entity.ShiftRequest) (*entity.ShiftRequest, error)
	Update(ctx context.Context, userID, resID, shiftID, shiftRequestID string, updateData map[string]interface{}) (*entity.ShiftRequest, error)
	Delete(ctx context.Context, userID, resID, shiftID, shiftRequestID string) error
	FindByID(ctx context.Context, userID, resID, shiftID, shiftRequestID string) (*entity.ShiftRequest, error)
	FindAllByShiftID(ctx context.Context, userID, resID, shiftID string) ([]*entity.ShiftRequest, error)
	FindAllByUserID(ctx context.Context, userID, resID, targetUserID string) ([]*entity.ShiftRequest, error)
}

type shiftRequestUseCase struct {
	shiftRequestRepo   repository.ShiftRequestRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewShiftRequestUseCase(shiftRequestRepo repository.ShiftRequestRepository, userRestaurantRepo repository.UserRestaurantRepository) ShiftRequestUseCase {
	return &shiftRequestUseCase{
		shiftRequestRepo:   shiftRequestRepo,
		userRestaurantRepo: userRestaurantRepo,
	}
}

func (u *shiftRequestUseCase) Create(ctx context.Context, userID, resID, shiftID string, shiftRequest *entity.ShiftRequest) (*entity.ShiftRequest, error) {
	if userID != shiftRequest.UserID.String() {
		return nil, xerror.Forbidden("You are not allowed to create shift request for another user")
	}

	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift request")
	}

	newShiftRequest, err := u.shiftRequestRepo.Create(ctx, shiftRequest)

	if err != nil {
		return nil, xerror.Internal("Can not create shift request")
	}

	return newShiftRequest, nil
}

func (u *shiftRequestUseCase) Update(ctx context.Context, userID, resID, shiftID, shiftRequestID string, updateData map[string]interface{}) (*entity.ShiftRequest, error) {
	shiftRequest, err := u.shiftRequestRepo.GetByID(ctx, shiftID, shiftRequestID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift request is not found")
		}
		return nil, xerror.Internal("Database failed")
	}
	sanitizedUpdateData := make(map[string]interface{}, len(updateData))
	for key, value := range updateData {
		sanitizedUpdateData[key] = value
	}

	isOwner := userID == shiftRequest.UserID.String()
	_, updatesStatus := sanitizedUpdateData["status"]
	isManager := false

	if !isOwner || updatesStatus {
		var err error
		isManager, err = u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)
		if err != nil {
			return nil, xerror.Internal("Can not check authority")
		}
	}

	if isOwner && !isManager {
		delete(sanitizedUpdateData, "status")
	} else if !isOwner && !isManager {
		return nil, xerror.Forbidden("You are not allowed to update shift request")
	}

	updatedShiftRequest, err := u.shiftRequestRepo.Update(ctx, sanitizedUpdateData, shiftID, shiftRequestID)

	if err != nil {
		return nil, xerror.Internal("Can not update shift request")
	}

	return updatedShiftRequest, nil
}

func (u *shiftRequestUseCase) Delete(ctx context.Context, userID, resID, shiftID, shiftRequestID string) error {
	shiftRequest, err := u.shiftRequestRepo.GetByID(ctx, shiftID, shiftRequestID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("Shift request is not found")
		}
		return xerror.Internal("Database failed")
	}

	if userID != shiftRequest.UserID.String() {
		isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

		if err != nil {
			return xerror.Internal("Can not check authority")
		}

		if !isAuthority {
			return xerror.Forbidden("You are not allowed to delete shift request")
		}
	}

	if err := u.shiftRequestRepo.Delete(ctx, shiftID, shiftRequestID); err != nil {
		return xerror.Internal("Can not delete shift request")
	}
	return nil
}

func (u *shiftRequestUseCase) FindByID(ctx context.Context, userID, resID, shiftID, shiftRequestID string) (*entity.ShiftRequest, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view shift request")
	}

	shiftRequest, err := u.shiftRequestRepo.GetByID(ctx, shiftID, shiftRequestID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift request is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftRequest, nil
}

func (u *shiftRequestUseCase) FindAllByShiftID(ctx context.Context, userID, resID, shiftID string) ([]*entity.ShiftRequest, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view shift requests")
	}

	shiftRequests, err := u.shiftRequestRepo.GetAllByShiftID(ctx, shiftID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift requests are not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftRequests, nil
}

func (u *shiftRequestUseCase) FindAllByUserID(ctx context.Context, userID, resID, targetUserID string) ([]*entity.ShiftRequest, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view shift requests")
	}

	shiftRequests, err := u.shiftRequestRepo.GetAllByUserID(ctx, targetUserID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift requests are not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftRequests, nil
}
