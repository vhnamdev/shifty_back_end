package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapPositionModelToEntity(input *model.CreatePositionInput) (*entity.Position, error) {
	position := &entity.Position{}
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}
	if input.Name != "" {
		position.Name = input.Name
	}

	if input.Description != "" {
		position.Description = input.Description
	}

	if input.Rank != nil {
		position.Rank = *input.Rank
	}

	if input.RestaurantID != "" {
		position.RestaurantID = uuid.MustParse(input.RestaurantID)
	}

	if input.Salary != nil {
		val := int64(*input.Salary)
		position.Salary = &val
	}

	if input.CanDeleteRestaurant != nil {
		position.CanDeleteRestaurant = *input.CanDeleteRestaurant
	}

	if input.CanUpdateRestaurant != nil {
		position.CanUpdateRestaurant = *input.CanUpdateRestaurant
	}

	return position, nil
}

func MapPositionEntityToModel(position *entity.Position) (*model.Position, error) {
	if position == nil {
		return nil, xerror.BadRequest("Position is not valid")
	}
	var salary *int
	if position.Salary != nil {
		val := int(*position.Salary)
		salary = &val
	}

	return &model.Position{
		Name:                position.Name,
		Description:         position.Description,
		Salary:              salary,
		Rank:                position.Rank,
		CanUpdateRestaurant: position.CanUpdateRestaurant,
		CanDeleteRestaurant: position.CanDeleteRestaurant,
		RestaurantID:        position.RestaurantID.String(),
	}, nil
}

func MapUpdatePositionToEntity(input *model.UpdatePositionInput) (map[string]interface{}, error) {
	updateData := make(map[string]interface{})
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}
	var salary *int64
	if input.Name != nil {
		updateData["name"] = input.Name
	}

	if input.Description != nil {
		updateData["description"] = input.Description
	}

	if input.Rank != nil {
		updateData["rank"] = input.Rank
	}
	if input.Salary != nil {
		val := int64(*input.Salary)
		salary = &val
	}
	updateData["salary"] = salary

	if input.CanDeleteRestaurant != nil {
		updateData["can_update_restaurant"] = input.CanDeleteRestaurant
	}

	if input.CanUpdateRestaurant != nil {
		updateData["can_delete_restaurant"] = input.CanDeleteRestaurant
	}
	return updateData, nil
}
