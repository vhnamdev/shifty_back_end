package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/xerror"
)

func MapReactionEntityToModel(reaction *entity.Reaction) (*model.Reaction, error) {
	if reaction == nil {
		return nil, xerror.BadRequest("Reaction is not valid")
	}

	return &model.Reaction{
		ID:        reaction.ID.String(),
		Type:      model.ReactionType(reaction.Type),
		PostID:    reaction.PostID.String(),
		AuthorID:  reaction.AuthorID.String(),
		CreatedAt: reaction.CreatedAt,
		UpdatedAt: reaction.UpdatedAt,
	}, nil
}

func MapReactionSummaryToModel(summary *usecase.ReactionSummary) (*model.ReactionSummary, error) {
	if summary == nil {
		return nil, xerror.BadRequest("Reaction summary is not valid")
	}

	counts := make([]*model.ReactionCount, 0, len(constants.ReactionTypes))
	for _, reactionType := range constants.ReactionTypes {
		counts = append(counts, &model.ReactionCount{
			Type:  model.ReactionType(reactionType),
			Count: summary.Counts[reactionType],
		})
	}

	var myReaction *model.Reaction
	if summary.MyReaction != nil {
		mappedReaction, err := MapReactionEntityToModel(summary.MyReaction)
		if err != nil {
			return nil, err
		}
		myReaction = mappedReaction
	}

	return &model.ReactionSummary{
		PostID:     summary.PostID,
		Total:      summary.Total,
		Counts:     counts,
		MyReaction: myReaction,
	}, nil
}

func MapReactionResultToModel(result *usecase.ReactionResult) (*model.ReactionResult, error) {
	if result == nil {
		return nil, xerror.BadRequest("Reaction result is not valid")
	}

	summary, err := MapReactionSummaryToModel(result.Summary)
	if err != nil {
		return nil, err
	}

	var reaction *model.Reaction
	if result.Reaction != nil {
		mappedReaction, err := MapReactionEntityToModel(result.Reaction)
		if err != nil {
			return nil, err
		}
		reaction = mappedReaction
	}

	return &model.ReactionResult{
		Action:   model.ReactionAction(result.Action),
		Reaction: reaction,
		Summary:  summary,
	}, nil
}
