package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
)

func MapCreateCommentToEntity(input *model.CreateCommentInput) (*entity.Comment, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	postID, err := uuid.Parse(input.PostID)
	if err != nil {
		return nil, xerror.BadRequest("Invalid post ID")
	}

	comment := &entity.Comment{
		Content: input.Content,
		PostID:  postID,
	}

	if input.ImageURL != nil {
		comment.ImageUrl = input.ImageURL
	}

	if input.ParentID != nil {
		parentID, err := uuid.Parse(*input.ParentID)
		if err != nil {
			return nil, xerror.BadRequest("Invalid parent comment ID")
		}
		comment.ParentID = &parentID
	}

	return comment, nil
}

func MapUpdateCommentToEntity(input *model.UpdateCommentInput) (map[string]interface{}, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	updateData := make(map[string]interface{})

	if input.Content != nil {
		updateData["content"] = *input.Content
	}

	if input.ImageURL != nil {
		updateData["image_url"] = *input.ImageURL
	}

	return updateData, nil
}

func MapCommentEntityToModel(comment *entity.Comment) (*model.Comment, error) {
	if comment == nil {
		return nil, xerror.BadRequest("Comment is not valid")
	}

	var parentID *string
	if comment.ParentID != nil {
		parentIDValue := comment.ParentID.String()
		parentID = &parentIDValue
	}

	replies := make([]*model.Comment, 0, len(comment.Replies))
	for i := range comment.Replies {
		reply := comment.Replies[i]
		mappedReply, err := MapCommentEntityToModel(&reply)
		if err != nil {
			return nil, err
		}
		replies = append(replies, mappedReply)
	}

	return &model.Comment{
		ID:        comment.ID.String(),
		Content:   comment.Content,
		ImageURL:  comment.ImageUrl,
		PostID:    comment.PostID.String(),
		AuthorID:  comment.AuthorID.String(),
		ParentID:  parentID,
		Replies:   replies,
		IsDeleted: comment.IsDeleted,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		DeletedAt: comment.DeletedAt,
	}, nil
}

func MapCommentsToPagination(comments []*entity.Comment, total int64, page, limit int) (*model.CommentPagination, error) {
	commentModels := make([]*model.Comment, 0, len(comments))

	for _, comment := range comments {
		mappedComment, err := MapCommentEntityToModel(comment)
		if err != nil {
			return nil, err
		}
		commentModels = append(commentModels, mappedComment)
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &model.CommentPagination{
		Data:        commentModels,
		Total:       int(total),
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}
