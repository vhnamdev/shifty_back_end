package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"strings"

	"github.com/google/uuid"
)

type PostUseCase interface {
	Create(ctx context.Context, userID, resID string, post *entity.Post) (*entity.Post, error)
	Update(ctx context.Context, userID, resID, postID string, updateData map[string]interface{}) (*entity.Post, error)
	Delete(ctx context.Context, userID, resID, postID string) error
	FindByID(ctx context.Context, userID, resID, postID string) (*entity.Post, error)
	FindAllByRestaurantID(ctx context.Context, userID, resID string, page, limit int) ([]*entity.Post, int64, error)
}

type postUseCase struct {
	postRepo           repository.PostRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewPostUseCase(postRepo repository.PostRepository, userRestaurantRepo repository.UserRestaurantRepository) PostUseCase {
	return &postUseCase{
		postRepo:           postRepo,
		userRestaurantRepo: userRestaurantRepo,
	}
}

func (u *postUseCase) Create(ctx context.Context, userID, resID string, post *entity.Post) (*entity.Post, error) {
	if post == nil {
		return nil, xerror.BadRequest("Post is not valid")
	}

	if strings.TrimSpace(post.Content) == "" && strings.TrimSpace(post.ImageUrl) == "" {
		return nil, xerror.BadRequest("Post must have content or image")
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid user ID")
	}

	parsedResID, err := uuid.Parse(resID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid restaurant ID")
	}

	if post.AuthorID != uuid.Nil && post.AuthorID != parsedUserID {
		return nil, xerror.Forbidden("You are not allowed to create post for another user")
	}

	if post.RestaurantID != uuid.Nil && post.RestaurantID != parsedResID {
		return nil, xerror.BadRequest("Post restaurant does not match request restaurant")
	}

	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)
	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create post")
	}

	post.AuthorID = parsedUserID
	post.RestaurantID = parsedResID

	newPost, err := u.postRepo.Create(ctx, post)
	if err != nil {
		return nil, xerror.Internal("Can not create post")
	}

	return newPost, nil
}

func (u *postUseCase) Update(ctx context.Context, userID, resID, postID string, updateData map[string]interface{}) (*entity.Post, error) {
	post, err := u.postRepo.GetByID(ctx, postID, resID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Post is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	if userID != post.AuthorID.String() {
		return nil, xerror.Forbidden("You are not allowed to update post")
	}

	if len(updateData) == 0 {
		return post, nil
	}

	updatedPost, err := u.postRepo.Update(ctx, updateData, postID, resID)
	if err != nil {
		return nil, xerror.Internal("Can not update post")
	}

	return updatedPost, nil
}

func (u *postUseCase) Delete(ctx context.Context, userID, resID, postID string) error {
	post, err := u.postRepo.GetByID(ctx, postID, resID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("Post is not found")
		}
		return xerror.Internal("Database failed")
	}

	if userID != post.AuthorID.String() {
		isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)
		if err != nil {
			return xerror.Internal("Can not check authority")
		}

		if !isAuthority {
			return xerror.Forbidden("You are not allowed to delete post")
		}
	}

	if err := u.postRepo.Delete(ctx, postID, resID); err != nil {
		return xerror.Internal("Can not delete post")
	}
	return nil
}

func (u *postUseCase) FindByID(ctx context.Context, userID, resID, postID string) (*entity.Post, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)
	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view post")
	}

	post, err := u.postRepo.GetByID(ctx, postID, resID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Post is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return post, nil
}

func (u *postUseCase) FindAllByRestaurantID(ctx context.Context, userID, resID string, page, limit int) ([]*entity.Post, int64, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)
	if err != nil {
		return nil, 0, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, 0, xerror.Forbidden("You are not allowed to view posts")
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	posts, total, err := u.postRepo.GetAllByRestaurantID(ctx, resID, page, limit)
	if err != nil {
		return nil, 0, xerror.Internal("Database failed")
	}

	return posts, total, nil
}
