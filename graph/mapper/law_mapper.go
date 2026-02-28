package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"
)

func MapLawEntityToModel(law *entity.Law) (*model.Law, error) {
	if law == nil {
		return nil, xerror.BadRequest("Law is not valid")
	}

	return &model.Law{
		ID:            law.ID.String(),
		Name:          law.Name,
		Description:   law.Description,
		SeverityLevel: string(law.SeverityLevel),
		RestaurantID:  law.Restaurant.ID.String(),
	}, nil

}
