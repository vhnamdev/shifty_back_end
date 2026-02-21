package graph

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapUserUpdateInputToEntity(id string, input *model.UpdateUserInput) *entity.User {
	u := &entity.User{
		ID: uuid.MustParse(id),
	}

	if input.FullName != nil {
		u.FullName = *input.FullName
	}
	if input.Address != nil {
		u.Address = input.Address
	}
	if input.PhoneNumber != nil {
		u.PhoneNumber = input.PhoneNumber
	}
	return u
}

func MapStaffUpdateToMap(input *model.UpdateStaffByManagerInput) (map[string]interface{}, error) {
	updateData := make(map[string]interface{})

	if input.Position != nil {
		posID, err := uuid.Parse(*input.Position)
		if err != nil {
			return nil, xerror.BadRequest("Invalid Position ID format")
		}
		updateData["position_id"] = posID
	}

	if input.IsBanned != nil {
		updateData["is_banned"] = *input.IsBanned
	}

	return updateData, nil
}

func MapUserEntityToModel(u *entity.User) *model.User {
	if u == nil {
		return nil
	}
	jobs := make([]*model.UserJob, 0, len(u.UserRestaurants))
	for _, ur := range u.UserRestaurants {
		posName := ""
		if ur.Position.ID != uuid.Nil {
			posName = ur.Position.Name
		}

		resName := ""
		if ur.Restaurant.ID != uuid.Nil {
			resName = ur.Restaurant.Name
		}

		jobs = append(jobs, &model.UserJob{
			RestaurantID:   ur.RestaurantID.String(),
			RestaurantName: resName,
			Position:       posName,
		})
	}
	return &model.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		Avatar:      u.Avatar,
		Email:       u.Email,
		Role:        string(u.Role),
		Address:     u.Address,
		PhoneNumber: u.PhoneNumber,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		Jobs:        jobs,
	}
}

func MapUserRestaurantEntityToModel(ur *entity.UserRestaurant) *model.UserRestaurant {
	if ur == nil {
		return nil
	}
	position := ""
	if ur.Position.ID != uuid.Nil {
		position = ur.Position.Name
	} else {
		position = ur.PositionID.String()
	}

	restaurant := ""
	if ur.Restaurant.ID != uuid.Nil {
		restaurant = ur.Restaurant.Name
	} else {
		restaurant = ur.RestaurantID.String()
	}

	idStr := ur.UserID.String()

	return &model.UserRestaurant{
		UserID:     idStr,
		Restaurant: restaurant,
		Position:   position,
		IsBanned:   ur.IsBanned,
		JoinedAt:   ur.JoinedAt,
	}
}

func MapRestaurantModelToEntity(input *model.CreateRestaurantInput) *entity.Restaurant {
	restaurant := &entity.Restaurant{}

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

	return restaurant
}

func MapRestaurantEntityToModel(restaurant *entity.Restaurant) *model.Restaurant {
	if restaurant == nil {
		return nil
	}
	
	modelLaws := make([]*model.Law, 0, len(restaurant.Laws))

	for _, lawEntity := range restaurant.Laws {
		mappedLaw := MapLawEntityToModel(&lawEntity)
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
	}

}

func MapLawEntityToModel(law *entity.Law) *model.Law {
	if law == nil {
		return nil
	}

	return &model.Law{
		ID:            law.ID.String(),
		Name:          law.Name,
		Description:   law.Description,
		SeverityLevel: string(law.SeverityLevel),
		RestaurantID:  law.Restaurant.ID.String(),
	}

}

func MapRestaurantUpdateToMap(input *model.UpdateRestaurantInput) (map[string]interface{}, error) {
	updateData := make(map[string]interface{})

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
