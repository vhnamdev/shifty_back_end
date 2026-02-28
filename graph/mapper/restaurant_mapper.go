package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"
)

func MapRestaurantModelToEntity(input *model.CreateRestaurantInput) (*entity.Restaurant, error) {
	restaurant := &entity.Restaurant{}
	if input == nil {
		return nil, xerror.BadRequest("Input is required")
	}
	if input.Address != "" {
		restaurant.Address = input.Address
	}

	if input.Name != "" {
		restaurant.Name = input.Name
	}

	if input.PhoneNumber != "" {
		restaurant.Phone = input.PhoneNumber
	}

	if input.Email != "" {
		restaurant.Email = input.Email
	}

	return restaurant, nil
}

func MapRestaurantEntityToModel(restaurant *entity.Restaurant) (*model.Restaurant, error) {
	if restaurant == nil {
		return nil, xerror.BadRequest("restaurant is required")
	}

	modelLaws := make([]*model.Law, 0, len(restaurant.Laws))

	for _, lawEntity := range restaurant.Laws {
		mappedLaw, err := MapLawEntityToModel(&lawEntity)
		if err != nil {
			return nil, xerror.Internal("Can not map law from entity to model")
		}
		modelLaws = append(modelLaws, mappedLaw)
	}
	return &model.Restaurant{
		ID:          restaurant.ID.String(),
		Name:        restaurant.Name,
		Avatar:      restaurant.Avatar,
		Status:      restaurant.Status,
		Laws:        modelLaws,
		Address:     restaurant.Address,
		PhoneNumber: restaurant.Phone,
		Email:       restaurant.Email,
		CreatedAt:   restaurant.CreatedAt,
	}, nil

}

func MapRestaurantUpdateToMap(input *model.UpdateRestaurantInput) (map[string]interface{}, error) {
	updateData := make(map[string]interface{})
	if input == nil {
		return nil, xerror.BadRequest("Input is required")
	}
	if input.Name != nil {
		updateData["name"] = *input.Name
	}

	if input.Address != nil {
		updateData["address"] = *input.Address
	}

	if input.Email != nil {
		updateData["email"] = *input.Email
	}

	if input.PhoneNumber != nil {
		updateData["phone"] = *input.PhoneNumber
	}

	if input.Status != nil {
		updateData["status"] = *input.Status
	}

	return updateData, nil
}
