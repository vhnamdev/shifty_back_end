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
