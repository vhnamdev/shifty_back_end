package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/xerror"
)

type PositionUseCase interface {
	Create(ctx context.Context, position *entity.Position) (*entity.Position, error)
}

type positionUseCase struct {
	positionRepo repository.PositionRepository
}

func NewPositionUseCase(positionRepo repository.PositionRepository) PositionUseCase {
	return &positionUseCase{
		positionRepo: positionRepo,
	}
}

func (u *positionUseCase) Create(ctx context.Context, position *entity.Position) (*entity.Position, error) {
	result, err := u.positionRepo.Create(ctx, position)

	if err != nil {
		return nil, xerror.Internal("Can not create position")
	}

	return result, nil
}
