package mapper

import (
	"encoding/json"
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func MapCreateShiftRuleToEntity(input *model.CreateShiftRuleInput) (*entity.ShiftRule, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	restaurantID, err := uuid.Parse(input.ResID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid Restaurant ID format")
	}

	if !json.Valid([]byte(input.Config)) {
		return nil, xerror.BadRequest("Config must be a valid JSON string")
	}
	return &entity.ShiftRule{
		Type:         input.Type,
		Name:         input.Name,
		Config:       datatypes.JSON([]byte(input.Config)),
		RestaurantID: restaurantID,
	}, nil
}

func MapUpdateShiftRuleToEntity(input *model.UpdateShiftRuleInput) (map[string]interface{}, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	updateData := make(map[string]interface{})

	if input.Type != nil {
		updateData["type"] = *input.Type
	}

	if input.Name != nil {
		updateData["name"] = *input.Name
	}

	if input.Config != nil {
		if !json.Valid([]byte(*input.Config)) {
			return nil, xerror.BadRequest("Config must be a valid JSON string")
		}
		updateData["config"] = datatypes.JSON([]byte(*input.Config))
	}

	return updateData, nil
}

func MapShiftRuleEntityToModel(rule *entity.ShiftRule) (*model.ShiftRule, error) {
	if rule == nil {
		return nil, xerror.BadRequest("Shift rule is not valid")
	}

	return &model.ShiftRule{
		ID:           rule.ID.String(),
		Type:         rule.Type,
		Name:         rule.Name,
		Config:       string(rule.Config),
		IsDeleted:    rule.IsDeleted,
		RestaurantID: rule.RestaurantID.String(),
		CreatedAt:    rule.CreatedAt,
		UpdatedAt:    rule.UpdatedAt,
	}, nil
}
