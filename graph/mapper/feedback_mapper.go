package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"
)

func MapFeedbackEntityToModel(feedback *entity.Feedback) (*model.Feedback, error) {
	if feedback == nil {
		return nil, xerror.BadRequest("Feedback is not valid")
	}
	return &model.Feedback{ID: feedback.ID.String(), Content: feedback.Content, RestaurantID: feedback.RestaurantID.String(), MemberID: feedback.MemberID.String(), ReviewerID: feedback.ReviewerID.String(), IsDeleted: feedback.IsDeleted, CreatedAt: feedback.CreatedAt, UpdatedAt: feedback.UpdatedAt, DeletedAt: feedback.DeletedAt}, nil
}

func MapFeedbacksToPagination(feedbacks []*entity.Feedback, total int64, page, limit int) (*model.FeedbackPagination, error) {
	feedbackModels := make([]*model.Feedback, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		mappedFeedback, err := MapFeedbackEntityToModel(feedback)
		if err != nil {
			return nil, err
		}
		feedbackModels = append(feedbackModels, mappedFeedback)
	}
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return &model.FeedbackPagination{Data: feedbackModels, Total: int(total), CurrentPage: page, TotalPages: totalPages}, nil
}
