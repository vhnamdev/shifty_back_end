package unit_test

import (
	"context"
	"errors"
	"testing"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupCommentUseCase() (*MockCommentRepo, *MockPostRepo, *MockUserRestaurantRepo, usecase.CommentUseCase) {
	mockCommentRepo := new(MockCommentRepo)
	mockPostRepo := new(MockPostRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewCommentUseCase(mockCommentRepo, mockPostRepo, mockUserResRepo)

	return mockCommentRepo, mockPostRepo, mockUserResRepo, u
}

func TestCommentUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	resID := uuid.New()
	postID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		inputComment := &entity.Comment{Content: "Nice update", PostID: postID}
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		expectedComment := &entity.Comment{ID: uuid.New(), Content: "Nice update", PostID: postID, AuthorID: userID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockCommentRepo.On("Create", ctx, inputComment).Return(expectedComment, nil)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputComment)

		assert.NoError(t, err)
		assert.Equal(t, expectedComment, res)
		assert.Equal(t, userID, inputComment.AuthorID)
		assert.Equal(t, postID, inputComment.PostID)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Success Image Only", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		imageURL := "https://example.com/comment.png"
		inputComment := &entity.Comment{ImageUrl: &imageURL, PostID: postID}
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		expectedComment := &entity.Comment{ID: uuid.New(), ImageUrl: &imageURL, PostID: postID, AuthorID: userID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockCommentRepo.On("Create", ctx, inputComment).Return(expectedComment, nil)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputComment)

		assert.NoError(t, err)
		assert.Equal(t, expectedComment, res)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Success Reply", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		parentID := uuid.New()
		inputComment := &entity.Comment{Content: "Reply comment", PostID: postID, ParentID: &parentID}
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		parentComment := &entity.Comment{ID: parentID, Content: "Parent comment", PostID: postID}
		expectedComment := &entity.Comment{ID: uuid.New(), Content: "Reply comment", PostID: postID, ParentID: &parentID, AuthorID: userID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, parentID.String(), postID.String()).Return(parentComment, nil)
		mockCommentRepo.On("Create", ctx, inputComment).Return(expectedComment, nil)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputComment)

		assert.NoError(t, err)
		assert.Equal(t, expectedComment, res)
		assert.Equal(t, userID, inputComment.AuthorID)
		assert.Equal(t, postID, inputComment.PostID)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Fail Empty Comment", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		inputComment := &entity.Comment{Content: "   ", PostID: postID}

		res, err := u.Create(ctx, userID.String(), resID.String(), inputComment)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Comment must have content or image")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockCommentRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		inputComment := &entity.Comment{Content: "Nice update", PostID: postID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(false, nil)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputComment)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create comment")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockCommentRepo.AssertNotCalled(t, "Create")
		mockUserResRepo.AssertExpectations(t)
	})

	t.Run("Fail Post Not Found", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		inputComment := &entity.Comment{Content: "Nice update", PostID: postID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(nil, gorm.ErrRecordNotFound)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputComment)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Post is not found")
		mockCommentRepo.AssertNotCalled(t, "Create")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Fail Parent Not Found", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		parentID := uuid.New()
		inputComment := &entity.Comment{Content: "Reply comment", PostID: postID, ParentID: &parentID}
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, parentID.String(), postID.String()).Return(nil, gorm.ErrRecordNotFound)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputComment)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Parent comment is not found")
		mockCommentRepo.AssertNotCalled(t, "Create")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})
}

func TestCommentUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	resID := uuid.New().String()
	postID := uuid.New().String()
	commentID := uuid.New().String()
	updateData := map[string]interface{}{"content": "Updated comment"}

	t.Run("Success Owner", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: userID, PostID: uuid.MustParse(postID), Content: "Old comment"}
		expectedComment := &entity.Comment{ID: existingComment.ID, AuthorID: userID, PostID: existingComment.PostID, Content: "Updated comment"}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)
		mockCommentRepo.On("Update", ctx, updateData, commentID, postID).Return(expectedComment, nil)

		res, err := u.Update(ctx, userID.String(), resID, postID, commentID, updateData)

		assert.NoError(t, err)
		assert.Equal(t, expectedComment, res)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Success Empty Update", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: userID, PostID: uuid.MustParse(postID), Content: "Existing comment"}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)

		res, err := u.Update(ctx, userID.String(), resID, postID, commentID, map[string]interface{}{})

		assert.NoError(t, err)
		assert.Equal(t, existingComment, res)
		mockCommentRepo.AssertNotCalled(t, "Update")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: otherUserID, PostID: uuid.MustParse(postID)}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)

		res, err := u.Update(ctx, userID.String(), resID, postID, commentID, updateData)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to update comment")
		mockCommentRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Fail Empty Result", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		imageURL := "https://example.com/comment.png"
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: userID, PostID: uuid.MustParse(postID), ImageUrl: &imageURL}
		emptyImage := "   "
		emptyUpdate := map[string]interface{}{"image_url": emptyImage}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)

		res, err := u.Update(ctx, userID.String(), resID, postID, commentID, emptyUpdate)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Comment must have content or image")
		mockCommentRepo.AssertNotCalled(t, "Update")
	})
}

func TestCommentUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	resID := uuid.New().String()
	postID := uuid.New().String()
	commentID := uuid.New().String()

	t.Run("Success Owner", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: userID, PostID: uuid.MustParse(postID)}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)
		mockCommentRepo.On("Delete", ctx, commentID, postID).Return(nil)

		err := u.Delete(ctx, userID.String(), resID, postID, commentID)

		assert.NoError(t, err)
		mockUserResRepo.AssertNotCalled(t, "HasManagementAuthority")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Success Manager", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: otherUserID, PostID: uuid.MustParse(postID)}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(true, nil)
		mockCommentRepo.On("Delete", ctx, commentID, postID).Return(nil)

		err := u.Delete(ctx, userID.String(), resID, postID, commentID)

		assert.NoError(t, err)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: otherUserID, PostID: uuid.MustParse(postID)}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(false, nil)

		err := u.Delete(ctx, userID.String(), resID, postID, commentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You are not allowed to delete comment")
		mockCommentRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Fail Repo Error", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		existingComment := &entity.Comment{ID: uuid.MustParse(commentID), AuthorID: userID, PostID: uuid.MustParse(postID)}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(existingComment, nil)
		mockCommentRepo.On("Delete", ctx, commentID, postID).Return(errors.New("db error"))

		err := u.Delete(ctx, userID.String(), resID, postID, commentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Can not delete comment")
	})
}

func TestCommentUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	postID := uuid.New().String()
	commentID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		expectedComment := &entity.Comment{ID: uuid.MustParse(commentID), Content: "Comment detail", PostID: uuid.MustParse(postID)}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(expectedComment, nil)

		res, err := u.FindByID(ctx, userID, resID, postID, commentID)

		assert.NoError(t, err)
		assert.Equal(t, expectedComment, res)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(false, nil)

		res, err := u.FindByID(ctx, userID, resID, postID, commentID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to view comment")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockCommentRepo.AssertNotCalled(t, "GetByID")
		mockUserResRepo.AssertExpectations(t)
	})

	t.Run("Fail Not Found", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetByID", ctx, commentID, postID).Return(nil, gorm.ErrRecordNotFound)

		res, err := u.FindByID(ctx, userID, resID, postID, commentID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Comment is not found")
	})
}

func TestCommentUseCase_FindAllByPostID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	postID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		replyID := uuid.New()
		expectedComments := []*entity.Comment{
			{ID: uuid.New(), Content: "First comment", Replies: []entity.Comment{{ID: replyID, Content: "Reply comment", PostID: uuid.MustParse(postID)}}},
			{ID: uuid.New(), Content: "Second comment"},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetAllByPostID", ctx, postID, 2, 5).Return(expectedComments, int64(12), nil)

		comments, total, err := u.FindAllByPostID(ctx, userID, resID, postID, 2, 5)

		assert.NoError(t, err)
		assert.Equal(t, expectedComments, comments)
		assert.Equal(t, int64(12), total)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Success Default Pagination", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()
		existingPost := &entity.Post{ID: uuid.MustParse(postID)}
		expectedComments := []*entity.Comment{{ID: uuid.New(), Content: "Default page"}}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockCommentRepo.On("GetAllByPostID", ctx, postID, 1, 10).Return(expectedComments, int64(1), nil)

		comments, total, err := u.FindAllByPostID(ctx, userID, resID, postID, 0, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedComments, comments)
		assert.Equal(t, int64(1), total)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockCommentRepo, mockPostRepo, mockUserResRepo, u := setupCommentUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(false, nil)

		comments, total, err := u.FindAllByPostID(ctx, userID, resID, postID, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Equal(t, int64(0), total)
		assert.Contains(t, err.Error(), "You are not allowed to view comments")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockCommentRepo.AssertNotCalled(t, "GetAllByPostID")
		mockUserResRepo.AssertExpectations(t)
	})
}
