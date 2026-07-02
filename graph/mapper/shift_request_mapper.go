package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapCreateShiftRequestToEntity(input *model.CreateShiftRequestInput) (*entity.ShiftRequest, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	shiftID, err := uuid.Parse(input.ShiftID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid Shift ID format")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid User ID format")
	}

	shiftRequest := &entity.ShiftRequest{
		ShiftID:   shiftID,
		UserID:    userID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Note:      input.Note,
	}

	if input.PositionID != nil {
		posID, err := uuid.Parse(*input.PositionID)
		if err != nil {
			return nil, xerror.BadRequest("Invalid Position ID format")
		}
		shiftRequest.PositionID = &posID
	}

	return shiftRequest, nil
}

func MapUpdateShiftRequestToEntity(input *model.UpdateShiftRequestInput) (map[string]interface{}, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	updateData := make(map[string]interface{})

	if input.PositionID != nil {
		posID, err := uuid.Parse(*input.PositionID)
		if err != nil {
			return nil, xerror.BadRequest("Invalid Position ID format")
		}
		updateData["position_id"] = posID
	}

	if input.StartTime != nil {
		updateData["start_time"] = *input.StartTime
	}

	if input.EndTime != nil {
		updateData["end_time"] = *input.EndTime
	}

	if input.Status != nil {
		updateData["status"] = *input.Status
	}

	if input.Note != nil {
		updateData["note"] = input.Note
	}

	return updateData, nil
}

func MapShiftRequestEntityToModel(shiftRequest *entity.ShiftRequest) (*model.ShiftRequest, error) {
	if shiftRequest == nil {
		return nil, xerror.BadRequest("Shift request is not valid")
	}

	modelShiftRequest := &model.ShiftRequest{
		ID:        shiftRequest.ID.String(),
		StartTime: shiftRequest.StartTime,
		EndTime:   shiftRequest.EndTime,
		Status:    shiftRequest.Status,
		UserID:    shiftRequest.UserID.String(),
		ShiftID:   shiftRequest.ShiftID.String(),
		Note:      shiftRequest.Note,
		IsDeleted: shiftRequest.IsDeleted,
		CreatedAt: shiftRequest.CreatedAt,
		UpdatedAt: shiftRequest.UpdatedAt,
	}

	if shiftRequest.PositionID != nil {
		posIDStr := shiftRequest.PositionID.String()
		modelShiftRequest.PositionID = &posIDStr
	}

	return modelShiftRequest, nil
}
