package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type ScheduleUseCase interface {
	Create(ctx context.Context, userID, resID string, schedule *entity.Schedule) (*entity.Schedule, error)
	Update(ctx context.Context, userID, resID, scheID string, updateData map[string]interface{}) (*entity.Schedule, error)
	Delete(ctx context.Context, userID, resID, scheID string) error
	FindByID(ctx context.Context, userID, resID, scheID string) (*entity.Schedule, error)
	FindAllByResID(ctx context.Context, userID, resID string) ([]*entity.Schedule, error)
}

type scheduleUseCase struct {
	scheduleRepo       repository.ScheduleRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewScheduleUseCase(
	scheduleRepo repository.ScheduleRepository,
	userRestaurantRepo repository.UserRestaurantRepository) ScheduleUseCase {
	return &scheduleUseCase{
		scheduleRepo:       scheduleRepo,
		userRestaurantRepo: userRestaurantRepo,
	}
}

// Create schedule
func (u *scheduleUseCase) Create(ctx context.Context, userID, resID string, schedule *entity.Schedule) (*entity.Schedule, error) {

	// Check isAuthority if this user is not allow, return error
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create schedule")
	}

	// Create
	newSchedule, err := u.scheduleRepo.Create(ctx, schedule)

	if err != nil {
		return nil, xerror.Internal("Can not create schedule")
	}
	return newSchedule, nil
}

// Update schedule
func (u *scheduleUseCase) Update(ctx context.Context, userID, resID, scheID string, updateData map[string]interface{}) (*entity.Schedule, error) {

	// Same with create func, we have to check isAuthority
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to update schedule")
	}

	// Update
	updatedSchedule, err := u.scheduleRepo.Update(ctx, resID, scheID, updateData)

	if err != nil {
		return nil, xerror.Internal("Can not update schedule")
	}

	return updatedSchedule, nil
}

// Delete schedule
func (u *scheduleUseCase) Delete(ctx context.Context, userID, resID, scheID string) error {
	// Same with create and update func =))))), we have to check isAuthority
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return xerror.Forbidden("You are not allowed to delete schedule")
	}

	// Delete
	err = u.scheduleRepo.Delete(ctx, scheID, resID)

	if err != nil {
		return xerror.Internal("Can not delete schedule")
	}

	return nil
}

// Find by shedule's ID
func (u *scheduleUseCase) FindByID(ctx context.Context, userID, resID, scheID string) (*entity.Schedule, error) {

	// Hehe :) not the same with create, update or delete.
	// This func I want to check this user who is find this shedule is in this restaurant or not
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}
	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to find schedule in this restaurant")
	}

	schedule, err := u.scheduleRepo.FindByID(ctx, scheID, resID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Schedule is not found")
		}

		return nil, xerror.Internal("Database failed")
	}

	return schedule, nil
}

// Find all by restaurant's ID
func (u *scheduleUseCase) FindAllByResID(ctx context.Context, userID, resID string) ([]*entity.Schedule, error) {
	// Hehe :) same with find by ID
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}
	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to find schedule in this restaurant")
	}

	schedules, err := u.scheduleRepo.FindAllByResID(ctx, resID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Schedules are not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return schedules, nil
}
