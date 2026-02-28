package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapUserUpdateInputToEntity(id string, input *model.UpdateUserInput) (*entity.User, error) {
	if input != nil && id != "" {
		return nil, xerror.BadRequest("input and id is required")
	}
	var user entity.User

	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	if input.Address != nil {
		user.Address = input.Address
	}
	if input.PhoneNumber != nil {
		user.PhoneNumber = input.PhoneNumber
	}
	return &user, nil
}

func MapStaffUpdateToMap(input *model.UpdateStaffByManagerInput) (map[string]interface{}, error) {
	updateData := make(map[string]interface{})
	if input == nil {
		return nil, xerror.BadRequest("Input is required")
	}
	if input.Position != nil {
		posID, err := uuid.Parse(*input.Position)
		if err != nil {
			return nil, err
		}
		updateData["position_id"] = posID
	}

	if input.IsBanned != nil {
		updateData["is_banned"] = *input.IsBanned
	}

	return updateData, nil
}

func MapUserEntityToModel(user *entity.User) (*model.User, error) {
	if user == nil {
		return nil, xerror.BadRequest("user is required")
	}
	jobs := make([]*model.UserJob, 0, len(user.UserRestaurants))
	for _, ur := range user.UserRestaurants {
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
		ID:          user.ID.String(),
		FullName:    user.FullName,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Role:        string(user.Role),
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		Jobs:        jobs,
	}, nil
}
