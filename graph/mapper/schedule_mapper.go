package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapScheduleInputToEntity(input *model.CreateScheduleInput) (*entity.Schedule, error) {
	var schedule entity.Schedule
	if input == nil {
		return nil, xerror.BadRequest("Input is required")
	}

	if input.RestaurantID != "" {
		schedule.RestaurantID = uuid.MustParse(input.RestaurantID)
	}

	if !input.StartTime.IsZero() {
		schedule.StartTime = input.StartTime
	}

	if !input.EndTime.IsZero() {
		schedule.EndTime = input.EndTime
	}

	return &schedule, nil
}

func MapScheduleEntityToModel(schedule *entity.Schedule) (*model.Schedule, error) {
	if schedule == nil {
		return nil, xerror.BadRequest("Schedule is required")
	}

	return &model.Schedule{
		RestaurantID:    schedule.RestaurantID.String(),
		StartTime:       schedule.StartTime,
		EndTime:         schedule.EndTime,
		CreatedAt:       schedule.CreatedAt,
		NumberOfMembers: schedule.NumberOfMembers,
		NumberOfShifts:  schedule.NumberOfShifts,
	}, nil
}

func MapScheduleUpdateToMap(input *model.UpdateScheduleInput) (map[string]interface{}, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is required")
	}

	updateData := make(map[string]interface{})

	if !input.EndTime.IsZero() {
		updateData["end_time"] = input.EndTime
	}

	if !input.StartTime.IsZero() {
		updateData["start_time"] = input.StartTime
	}

	return updateData, nil

}
