package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapUserRestaurantEntityToModel(ur *entity.UserRestaurant) (*model.UserRestaurant, error) {
	if ur == nil {
		return nil, xerror.BadRequest("User restaurant is not valid")
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
	}, nil
}
