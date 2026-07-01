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

type CommentUseCase interface {
	Create(ctx context.Context, userID, resID string, comment *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, userID, resID, postID, commentID string, updateData map[string]interface{}) (*entity.Comment, error)
	Delete(ctx context.Context, userID, resID, postID, commentID string) error
	FindByID(ctx context.Context, userID, resID, postID, commentID string) (*entity.Comment, error)
	FindAllByPostID(ctx context.Context, userID, resID, postID string, page, limit int) ([]*entity.Comment, int64, error)
}

type commentUseCase struct {
	commentRepo        repository.CommentRepository
	postRepo           repository.PostRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewCommentUseCase(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRestaurantRepo repository.UserRestaurantRepository) CommentUseCase {
	return &commentUseCase{
		commentRepo:        commentRepo,
		postRepo:           postRepo,
		userRestaurantRepo: userRestaurantRepo,
	}
}

func (u *commentUseCase) Create(ctx context.Context, userID, resID string, comment *entity.Comment) (*entity.Comment, error) {
	if comment == nil {
		return nil, xerror.BadRequest("Comment is not valid")
	}

	if strings.TrimSpace(comment.Content) == "" && isEmptyCommentImage(comment.ImageUrl) {
		return nil, xerror.BadRequest("Comment must have content or image")
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid user ID")
	}

	if comment.AuthorID != uuid.Nil && comment.AuthorID != parsedUserID {
		return nil, xerror.Forbidden("You are not allowed to create comment for another user")
	}

	post, err := u.getAccessiblePost(ctx, userID, resID, comment.PostID.String(), "create comment")
	if err != nil {
		return nil, err
	}

	if comment.ParentID != nil {
		if _, err := u.commentRepo.GetByID(ctx, comment.ParentID.String(), post.ID.String()); err != nil {
			if utils.IsRecordNotFoundError(err) {
				return nil, xerror.NotFound("Parent comment is not found")
			}
			return nil, xerror.Internal("Database failed")
		}
	}

	comment.AuthorID = parsedUserID
	comment.PostID = post.ID

	newComment, err := u.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, xerror.Internal("Can not create comment")
	}

	return newComment, nil
}

func (u *commentUseCase) Update(ctx context.Context, userID, resID, postID, commentID string, updateData map[string]interface{}) (*entity.Comment, error) {
	if _, err := u.getAccessiblePost(ctx, userID, resID, postID, "update comment"); err != nil {
		return nil, err
	}

	comment, err := u.commentRepo.GetByID(ctx, commentID, postID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Comment is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	if userID != comment.AuthorID.String() {
		return nil, xerror.Forbidden("You are not allowed to update comment")
	}

	if len(updateData) == 0 {
		return comment, nil
	}

	nextContent := comment.Content
	if content, ok := updateData["content"].(string); ok {
		nextContent = content
	}

	nextImageURL := comment.ImageUrl
	if imageURL, ok := updateData["image_url"].(string); ok {
		nextImageURL = &imageURL
	}

	if strings.TrimSpace(nextContent) == "" && isEmptyCommentImage(nextImageURL) {
		return nil, xerror.BadRequest("Comment must have content or image")
	}

	updatedComment, err := u.commentRepo.Update(ctx, updateData, commentID, postID)
	if err != nil {
		return nil, xerror.Internal("Can not update comment")
	}

	return updatedComment, nil
}

func (u *commentUseCase) Delete(ctx context.Context, userID, resID, postID, commentID string) error {
	if _, err := u.getAccessiblePost(ctx, userID, resID, postID, "delete comment"); err != nil {
		return err
	}

	comment, err := u.commentRepo.GetByID(ctx, commentID, postID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("Comment is not found")
		}
		return xerror.Internal("Database failed")
	}

	if userID != comment.AuthorID.String() {
		isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)
		if err != nil {
			return xerror.Internal("Can not check authority")
		}

		if !isAuthority {
			return xerror.Forbidden("You are not allowed to delete comment")
		}
	}

	if err := u.commentRepo.Delete(ctx, commentID, postID); err != nil {
		return xerror.Internal("Can not delete comment")
	}
	return nil
}

func (u *commentUseCase) FindByID(ctx context.Context, userID, resID, postID, commentID string) (*entity.Comment, error) {
	if _, err := u.getAccessiblePost(ctx, userID, resID, postID, "view comment"); err != nil {
		return nil, err
	}

	comment, err := u.commentRepo.GetByID(ctx, commentID, postID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Comment is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return comment, nil
}

func (u *commentUseCase) FindAllByPostID(ctx context.Context, userID, resID, postID string, page, limit int) ([]*entity.Comment, int64, error) {
	if _, err := u.getAccessiblePost(ctx, userID, resID, postID, "view comments"); err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	comments, total, err := u.commentRepo.GetAllByPostID(ctx, postID, page, limit)
	if err != nil {
		return nil, 0, xerror.Internal("Database failed")
	}

	return comments, total, nil
}

func (u *commentUseCase) getAccessiblePost(ctx context.Context, userID, resID, postID, action string) (*entity.Post, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)
	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to " + action)
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

func isEmptyCommentImage(imageURL *string) bool {
	return imageURL == nil || strings.TrimSpace(*imageURL) == ""
}
