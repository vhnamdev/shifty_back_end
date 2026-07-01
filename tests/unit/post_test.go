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

func setupPostUseCase() (*MockPostRepo, *MockUserRestaurantRepo, usecase.PostUseCase) {
	mockPostRepo := new(MockPostRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewPostUseCase(mockPostRepo, mockUserResRepo)

	return mockPostRepo, mockUserResRepo, u
}

func TestPostUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	resID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		inputPost := &entity.Post{Content: "New announcement"}
		expectedPost := &entity.Post{ID: uuid.New(), Content: "New announcement", AuthorID: userID, RestaurantID: resID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("Create", ctx, inputPost).Return(expectedPost, nil)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputPost)

		assert.NoError(t, err)
		assert.Equal(t, expectedPost, res)
		assert.Equal(t, userID, inputPost.AuthorID)
		assert.Equal(t, resID, inputPost.RestaurantID)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Success Image Only", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		inputPost := &entity.Post{ImageUrl: "https://example.com/image.png"}
		expectedPost := &entity.Post{ID: uuid.New(), ImageUrl: inputPost.ImageUrl, AuthorID: userID, RestaurantID: resID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("Create", ctx, inputPost).Return(expectedPost, nil)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputPost)

		assert.NoError(t, err)
		assert.Equal(t, expectedPost, res)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Fail Empty Post", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		inputPost := &entity.Post{Content: "   "}

		res, err := u.Create(ctx, userID.String(), resID.String(), inputPost)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Post must have content or image")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockPostRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Fail Different Author", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		inputPost := &entity.Post{Content: "New announcement", AuthorID: uuid.New()}

		res, err := u.Create(ctx, userID.String(), resID.String(), inputPost)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create post for another user")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockPostRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		inputPost := &entity.Post{Content: "New announcement"}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(false, nil)

		res, err := u.Create(ctx, userID.String(), resID.String(), inputPost)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create post")
		mockPostRepo.AssertNotCalled(t, "Create")
		mockUserResRepo.AssertExpectations(t)
	})
}

func TestPostUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	resID := uuid.New().String()
	postID := uuid.New().String()
	updateData := map[string]interface{}{"content": "Updated post"}

	t.Run("Success Owner", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		existingPost := &entity.Post{ID: uuid.New(), AuthorID: userID}
		expectedPost := &entity.Post{ID: existingPost.ID, AuthorID: userID, Content: "Updated post"}

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockPostRepo.On("Update", ctx, updateData, postID, resID).Return(expectedPost, nil)

		res, err := u.Update(ctx, userID.String(), resID, postID, updateData)

		assert.NoError(t, err)
		assert.Equal(t, expectedPost, res)
		mockUserResRepo.AssertNotCalled(t, "HasManagementAuthority")
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Success Empty Update", func(t *testing.T) {
		mockPostRepo, _, u := setupPostUseCase()
		existingPost := &entity.Post{ID: uuid.New(), AuthorID: userID, Content: "Existing post"}

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)

		res, err := u.Update(ctx, userID.String(), resID, postID, map[string]interface{}{})

		assert.NoError(t, err)
		assert.Equal(t, existingPost, res)
		mockPostRepo.AssertNotCalled(t, "Update")
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockPostRepo, _, u := setupPostUseCase()
		existingPost := &entity.Post{ID: uuid.New(), AuthorID: otherUserID}

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)

		res, err := u.Update(ctx, userID.String(), resID, postID, updateData)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to update post")
		mockPostRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Fail Not Found", func(t *testing.T) {
		mockPostRepo, _, u := setupPostUseCase()

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(nil, gorm.ErrRecordNotFound)

		res, err := u.Update(ctx, userID.String(), resID, postID, updateData)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Post is not found")
		mockPostRepo.AssertNotCalled(t, "Update")
	})
}

func TestPostUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	resID := uuid.New().String()
	postID := uuid.New().String()

	t.Run("Success Owner", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		existingPost := &entity.Post{ID: uuid.New(), AuthorID: userID}

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockPostRepo.On("Delete", ctx, postID, resID).Return(nil)

		err := u.Delete(ctx, userID.String(), resID, postID)

		assert.NoError(t, err)
		mockUserResRepo.AssertNotCalled(t, "HasManagementAuthority")
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Success Manager", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		existingPost := &entity.Post{ID: uuid.New(), AuthorID: otherUserID}

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(true, nil)
		mockPostRepo.On("Delete", ctx, postID, resID).Return(nil)

		err := u.Delete(ctx, userID.String(), resID, postID)

		assert.NoError(t, err)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		existingPost := &entity.Post{ID: uuid.New(), AuthorID: otherUserID}

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(false, nil)

		err := u.Delete(ctx, userID.String(), resID, postID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You are not allowed to delete post")
		mockPostRepo.AssertNotCalled(t, "Delete")
		mockUserResRepo.AssertExpectations(t)
	})

	t.Run("Fail Repo Error", func(t *testing.T) {
		mockPostRepo, _, u := setupPostUseCase()
		existingPost := &entity.Post{ID: uuid.New(), AuthorID: userID}

		mockPostRepo.On("GetByID", ctx, postID, resID).Return(existingPost, nil)
		mockPostRepo.On("Delete", ctx, postID, resID).Return(errors.New("db error"))

		err := u.Delete(ctx, userID.String(), resID, postID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Can not delete post")
	})
}

func TestPostUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	postID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		expectedPost := &entity.Post{ID: uuid.New(), Content: "Post detail"}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(expectedPost, nil)

		res, err := u.FindByID(ctx, userID, resID, postID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPost, res)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(false, nil)

		res, err := u.FindByID(ctx, userID, resID, postID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to view post")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockUserResRepo.AssertExpectations(t)
	})

	t.Run("Fail Not Found", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID, resID).Return(nil, gorm.ErrRecordNotFound)

		res, err := u.FindByID(ctx, userID, resID, postID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Post is not found")
	})
}

func TestPostUseCase_FindAllByRestaurantID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		expectedPosts := []*entity.Post{
			{ID: uuid.New(), Content: "First post"},
			{ID: uuid.New(), Content: "Second post"},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetAllByRestaurantID", ctx, resID, 2, 5).Return(expectedPosts, int64(12), nil)

		posts, total, err := u.FindAllByRestaurantID(ctx, userID, resID, 2, 5)

		assert.NoError(t, err)
		assert.Equal(t, expectedPosts, posts)
		assert.Equal(t, int64(12), total)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Success Default Pagination", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()
		expectedPosts := []*entity.Post{{ID: uuid.New(), Content: "Default page"}}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPostRepo.On("GetAllByRestaurantID", ctx, resID, 1, 10).Return(expectedPosts, int64(1), nil)

		posts, total, err := u.FindAllByRestaurantID(ctx, userID, resID, 0, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedPosts, posts)
		assert.Equal(t, int64(1), total)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockPostRepo, mockUserResRepo, u := setupPostUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(false, nil)

		posts, total, err := u.FindAllByRestaurantID(ctx, userID, resID, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Equal(t, int64(0), total)
		assert.Contains(t, err.Error(), "You are not allowed to view posts")
		mockPostRepo.AssertNotCalled(t, "GetAllByRestaurantID")
	})
}
