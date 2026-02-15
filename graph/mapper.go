package graph

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"

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

func MapUserEntityToModel(u *entity.User) *model.User {
	if u == nil {
		return nil
	}
	positionName := ""
	if u.Position != nil {
		positionName = u.Position.Name
	}
	restaurantName := ""
	if u.Restaurant != nil {
		restaurantName = u.Restaurant.Name
	}
	return &model.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		Avatar:      u.Avatar,
		Email:       u.Email,
		Role:        u.Role,
		Address:     u.Address,
		PhoneNumber: u.PhoneNumber,
		Position:    positionName,
		Restaurant:  restaurantName,
		CreatedAt:   u.CreatedAt,
	}
}
