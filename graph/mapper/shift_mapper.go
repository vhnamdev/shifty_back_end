package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapShiftModelToEntity(input *model.CreateShiftInput) (*entity.Shift, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	scheID, err := uuid.Parse(input.ScheduleID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid Schedule ID format")
	}

	shift := &entity.Shift{
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
		ScheduleID:      scheID,
		NumberOfMembers: input.NumberOfMembers,
		Type:            input.Type,
	}

	if input.IsHoliday != nil {
		shift.IsHoliday = *input.IsHoliday
	}

	if input.WageMultiplier != nil {
		shift.WageMultiplier = *input.WageMultiplier
	}

	return shift, nil
}

func MapShiftEntityToModel(shift *entity.Shift) (*model.Shift, error) {
	if shift == nil {
		return nil, xerror.BadRequest("Shift is not valid")
	}

	return &model.Shift{
		StartTime:  shift.StartTime,
		EndTime:    shift.EndTime,
		ScheduleID: shift.ScheduleID.String(),
		Type:       shift.Type,
		IsHoliday:  shift.IsHoliday,
	}, nil
}

func MapUpdateShiftToEntity(input *model.UpdateShiftInput) (map[string]interface{}, error) {
	updateData := make(map[string]interface{})

	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	if input.StartTime != nil {
        updateData["start_time"] = *input.StartTime
    }

    if input.EndTime != nil {
        updateData["end_time"] = *input.EndTime
    }

    if input.NumberOfMembers != nil {
        updateData["number_of_members"] = *input.NumberOfMembers
    }

    if input.Type != nil {
        updateData["type"] = *input.Type
    }

    if input.IsHoliday != nil {
        updateData["is_holiday"] = *input.IsHoliday
    }

    if input.WageMultiplier != nil {
        updateData["wage_multiplier"] = *input.WageMultiplier
    }

	return updateData,nil

}
