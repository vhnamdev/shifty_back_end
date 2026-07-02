package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapCreateShiftAssignmentToEntity(input *model.CreateShiftAssignmentInput) (*entity.ShiftAssignment, error) {
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

	assignment := &entity.ShiftAssignment{
		ShiftID:      shiftID,
		UserID:       userID,
		CheckInTime:  input.CheckInTime,
		CheckOutTime: input.CheckOutTime,
		Note:         input.Note,
	}

	if input.PositionID != nil {
		posID, err := uuid.Parse(*input.PositionID)
		if err != nil {
			return nil, xerror.BadRequest("Invalid Position ID format")
		}
		assignment.PositionID = &posID
	}

	return assignment, nil
}

func MapUpdateShiftAssignmentToEntity(input *model.UpdateShiftAssignmentInput) (map[string]interface{}, error) {
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

	if input.CheckInTime != nil {
		updateData["check_in_time"] = *input.CheckInTime
	}

	if input.CheckOutTime != nil {
		updateData["check_out_time"] = *input.CheckOutTime
	}

	if input.Note != nil {
		updateData["note"] = input.Note
	}

	return updateData, nil
}

func MapShiftAssignmentEntityToModel(assignment *entity.ShiftAssignment) (*model.ShiftAssignment, error) {
	if assignment == nil {
		return nil, xerror.BadRequest("Shift assignment is not valid")
	}

	modelAssignment := &model.ShiftAssignment{
		ID:           assignment.ID.String(),
		ShiftID:      assignment.ShiftID.String(),
		UserID:       assignment.UserID.String(),
		CheckInTime:  assignment.CheckInTime,
		CheckOutTime: assignment.CheckOutTime,
		Note:         assignment.Note,
		IsDeleted:    assignment.IsDeleted,
		CreatedAt:    assignment.CreatedAt,
		UpdatedAt:    assignment.UpdatedAt,
	}

	if assignment.PositionID != nil {
		posIDStr := assignment.PositionID.String()
		modelAssignment.PositionID = &posIDStr
	}

	return modelAssignment, nil
}
