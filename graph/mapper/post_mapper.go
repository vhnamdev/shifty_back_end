package mapper

import (
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/xerror"
)

func MapCreatePostToEntity(input *model.CreatePostInput) (*entity.Post, error) {
	if input == nil {
		return nil, xerror.BadRequest("Input is not valid")
	}

	post := &entity.Post{
		Content: input.Content,
	}

	if input.ImageURL != nil {
		post.ImageUrl = *input.ImageURL
	}

	return post, nil
}

func MapUpdatePostToEntity(input *model.UpdatePostInput) (map[string]interface{}, error) {
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

func MapPostEntityToModel(post *entity.Post) (*model.Post, error) {
	if post == nil {
		return nil, xerror.BadRequest("Post is not valid")
	}

	return &model.Post{
		ID:           post.ID.String(),
		Content:      post.Content,
		ImageURL:     post.ImageUrl,
		RestaurantID: post.RestaurantID.String(),
		AuthorID:     post.AuthorID.String(),
		IsDeleted:    post.IsDeleted,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		DeletedAt:    post.DeletedAt,
	}, nil
}

func MapPostsToPagination(posts []*entity.Post, total int64, page, limit int) (*model.PostPagination, error) {
	postModels := make([]*model.Post, 0, len(posts))

	for _, post := range posts {
		mappedPost, err := MapPostEntityToModel(post)
		if err != nil {
			return nil, err
		}
		postModels = append(postModels, mappedPost)
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &model.PostPagination{
		Data:        postModels,
		Total:       int(total),
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}
