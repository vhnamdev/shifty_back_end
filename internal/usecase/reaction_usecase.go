package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"strings"

	"github.com/google/uuid"
)

const (
	ReactionActionCreated = "CREATED"
	ReactionActionUpdated = "UPDATED"
	ReactionActionDeleted = "DELETED"
)

type ReactionSummary struct {
	PostID     string
	Counts     map[string]int
	Total      int
	MyReaction *entity.Reaction
}

type ReactionResult struct {
	Action   string
	Reaction *entity.Reaction
	Summary  *ReactionSummary
}

type ReactionUseCase interface {
	ReactToPost(ctx context.Context, userID, resID, postID, reactionType string) (*ReactionResult, error)
	FindSummaryByPostID(ctx context.Context, userID, resID, postID string) (*ReactionSummary, error)
}

type reactionUseCase struct {
	reactionRepo       repository.ReactionRepository
	postRepo           repository.PostRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewReactionUseCase(reactionRepo repository.ReactionRepository, postRepo repository.PostRepository, userRestaurantRepo repository.UserRestaurantRepository) ReactionUseCase {
	return &reactionUseCase{
		reactionRepo:       reactionRepo,
		postRepo:           postRepo,
		userRestaurantRepo: userRestaurantRepo,
	}
}

func (u *reactionUseCase) ReactToPost(ctx context.Context, userID, resID, postID, reactionType string) (*ReactionResult, error) {
	normalizedType := strings.ToUpper(strings.TrimSpace(reactionType))
	if !constants.IsValidReactionType(normalizedType) {
		return nil, xerror.BadRequest("Reaction type is not valid")
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid user ID")
	}

	post, err := u.getAccessiblePost(ctx, userID, resID, postID, "react to post")
	if err != nil {
		return nil, err
	}

	existingReaction, err := u.reactionRepo.GetByPostAndAuthor(ctx, post.ID.String(), userID)
	if err != nil && !utils.IsRecordNotFoundError(err) {
		return nil, xerror.Internal("Database failed")
	}

	var action string
	var reaction *entity.Reaction

	if utils.IsRecordNotFoundError(err) {
		reaction, err = u.reactionRepo.Create(ctx, &entity.Reaction{
			Type:     normalizedType,
			PostID:   post.ID,
			AuthorID: parsedUserID,
		})
		if err != nil {
			return nil, xerror.Internal("Can not create reaction")
		}
		action = ReactionActionCreated
	} else if existingReaction.Type == normalizedType {
		if err := u.reactionRepo.Delete(ctx, existingReaction.ID.String()); err != nil {
			return nil, xerror.Internal("Can not delete reaction")
		}
		action = ReactionActionDeleted
	} else {
		reaction, err = u.reactionRepo.UpdateType(ctx, existingReaction.ID.String(), normalizedType)
		if err != nil {
			return nil, xerror.Internal("Can not update reaction")
		}
		action = ReactionActionUpdated
	}

	summary, err := u.buildSummary(ctx, userID, post.ID.String())
	if err != nil {
		return nil, err
	}

	return &ReactionResult{
		Action:   action,
		Reaction: reaction,
		Summary:  summary,
	}, nil
}

func (u *reactionUseCase) FindSummaryByPostID(ctx context.Context, userID, resID, postID string) (*ReactionSummary, error) {
	post, err := u.getAccessiblePost(ctx, userID, resID, postID, "view reactions")
	if err != nil {
		return nil, err
	}

	return u.buildSummary(ctx, userID, post.ID.String())
}

func (u *reactionUseCase) getAccessiblePost(ctx context.Context, userID, resID, postID, action string) (*entity.Post, error) {
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

func (u *reactionUseCase) buildSummary(ctx context.Context, userID, postID string) (*ReactionSummary, error) {
	counts, total, err := u.reactionRepo.CountByPostID(ctx, postID)
	if err != nil {
		return nil, xerror.Internal("Database failed")
	}

	myReaction, err := u.reactionRepo.GetByPostAndAuthor(ctx, postID, userID)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			myReaction = nil
		} else {
			return nil, xerror.Internal("Database failed")
		}
	}

	return &ReactionSummary{
		PostID:     postID,
		Counts:     counts,
		Total:      total,
		MyReaction: myReaction,
	}, nil
}
