package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
)

type UserUseCase interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByID(ctx context.Context, ID string) (*entity.User, error)
}
type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

// Find user by email
func (u *userUseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	// Get user
	user, err := u.userRepo.GetByEmail(ctx, email)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("User not found")
		}
		return nil, xerror.Internal("Database error")
	}

	return user, nil
}

// Find user by ID
func (u *userUseCase) FindUserByID(ctx context.Context, ID string) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, ID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("User not found")
		}
		return nil, xerror.Internal("Database error")
	}

	return user, nil
}
