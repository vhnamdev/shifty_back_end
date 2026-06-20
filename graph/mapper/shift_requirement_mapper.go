package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapCreateRequirementToEntity(input *model.CreateRequirementInput) (*entity.ShiftRequirement, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	shiftID, err := uuid.Parse(input.ShiftID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid Shift ID format")
	}

	req := &entity.ShiftRequirement{
		ShiftID:   shiftID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Quantity:  input.Quantity,
		Note:      input.Note,
	}


	if input.PositionID != nil {
		posID, err := uuid.Parse(*input.PositionID)
		if err != nil {
			return nil, xerror.BadRequest("Invalid Position ID format")
		}
		req.PositionID = &posID
	}

	return req, nil
}


func MapUpdateRequirementToEntity(input *model.UpdateRequirementInput) (map[string]interface{}, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	updateData := make(map[string]interface{})

	if input.StartTime != nil {
		updateData["start_time"] = *input.StartTime
	}

	if input.EndTime != nil {
		updateData["end_time"] = *input.EndTime
	}

	if input.Quantity != nil {
		updateData["quantity"] = *input.Quantity
	}

	if input.Note != nil {
		updateData["note"] = input.Note
	}

	if input.PositionID != nil {
		posID, err := uuid.Parse(*input.PositionID)
		if err != nil {
			return nil, xerror.BadRequest("Invalid Position ID format")
		}
		updateData["position_id"] = posID
	}

	return updateData, nil
}


func MapRequirementEntityToModel(req *entity.ShiftRequirement) (*model.ShiftRequirement, error) {
	if req == nil {
		return nil, xerror.BadRequest("Requirement is not valid")
	}

	modelReq := &model.ShiftRequirement{
		ID:        req.ID.String(),
		ShiftID:   req.ShiftID.String(),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Quantity:  req.Quantity,
		Note:      req.Note,
		IsDeleted: req.IsDeleted,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}

	if req.PositionID != nil {
		posIDStr := req.PositionID.String()
		modelReq.PositionID = &posIDStr
	}

	return modelReq, nil
}