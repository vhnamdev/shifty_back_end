package unit_test

import (
	"context"
	"errors"
	"testing"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/constants"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupReactionUseCase() (*MockReactionRepo, *MockPostRepo, *MockUserRestaurantRepo, usecase.ReactionUseCase) {
	mockReactionRepo := new(MockReactionRepo)
	mockPostRepo := new(MockPostRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewReactionUseCase(mockReactionRepo, mockPostRepo, mockUserResRepo)

	return mockReactionRepo, mockPostRepo, mockUserResRepo, u
}

func TestReactionUseCase_ReactToPost(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	resID := uuid.New()
	postID := uuid.New()

	t.Run("Success Create", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		expectedReaction := &entity.Reaction{ID: uuid.New(), Type: constants.ReactionLike, PostID: postID, AuthorID: userID}
		counts := map[string]int{constants.ReactionLike: 1}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(nil, gorm.ErrRecordNotFound).Once()
		mockReactionRepo.On("Create", ctx, mock.MatchedBy(func(reaction *entity.Reaction) bool {
			return reaction.Type == constants.ReactionLike &&
				reaction.PostID == postID &&
				reaction.AuthorID == userID
		})).Return(expectedReaction, nil)
		mockReactionRepo.On("CountByPostID", ctx, postID.String()).Return(counts, 1, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(expectedReaction, nil).Once()

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), " like ")

		assert.NoError(t, err)
		assert.Equal(t, usecase.ReactionActionCreated, res.Action)
		assert.Equal(t, expectedReaction, res.Reaction)
		assert.Equal(t, postID.String(), res.Summary.PostID)
		assert.Equal(t, 1, res.Summary.Total)
		assert.Equal(t, counts, res.Summary.Counts)
		assert.Equal(t, expectedReaction, res.Summary.MyReaction)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})

	t.Run("Success Delete Same Type", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		existingReaction := &entity.Reaction{ID: uuid.New(), Type: constants.ReactionLove, PostID: postID, AuthorID: userID}
		counts := map[string]int{}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(existingReaction, nil).Once()
		mockReactionRepo.On("Delete", ctx, existingReaction.ID.String()).Return(nil)
		mockReactionRepo.On("CountByPostID", ctx, postID.String()).Return(counts, 0, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), constants.ReactionLove)

		assert.NoError(t, err)
		assert.Equal(t, usecase.ReactionActionDeleted, res.Action)
		assert.Nil(t, res.Reaction)
		assert.Equal(t, 0, res.Summary.Total)
		assert.Nil(t, res.Summary.MyReaction)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})

	t.Run("Success Update Different Type", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		existingReaction := &entity.Reaction{ID: uuid.New(), Type: constants.ReactionLike, PostID: postID, AuthorID: userID}
		expectedReaction := &entity.Reaction{ID: existingReaction.ID, Type: constants.ReactionHaha, PostID: postID, AuthorID: userID}
		counts := map[string]int{constants.ReactionHaha: 1}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(existingReaction, nil).Once()
		mockReactionRepo.On("UpdateType", ctx, existingReaction.ID.String(), constants.ReactionHaha).Return(expectedReaction, nil)
		mockReactionRepo.On("CountByPostID", ctx, postID.String()).Return(counts, 1, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(expectedReaction, nil).Once()

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), constants.ReactionHaha)

		assert.NoError(t, err)
		assert.Equal(t, usecase.ReactionActionUpdated, res.Action)
		assert.Equal(t, expectedReaction, res.Reaction)
		assert.Equal(t, 1, res.Summary.Total)
		assert.Equal(t, expectedReaction, res.Summary.MyReaction)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})

	t.Run("Fail Invalid Type", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), "INVALID")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Reaction type is not valid")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockReactionRepo.AssertNotCalled(t, "GetByPostAndAuthor")
	})

	t.Run("Fail Invalid User ID", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()

		res, err := u.ReactToPost(ctx, "bad-user-id", resID.String(), postID.String(), constants.ReactionLike)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Invalid user ID")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockReactionRepo.AssertNotCalled(t, "GetByPostAndAuthor")
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(false, nil)

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), constants.ReactionLike)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to react to post")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockReactionRepo.AssertNotCalled(t, "GetByPostAndAuthor")
		mockUserResRepo.AssertExpectations(t)
	})

	t.Run("Fail Post Not Found", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(nil, gorm.ErrRecordNotFound)

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), constants.ReactionLike)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Post is not found")
		mockReactionRepo.AssertNotCalled(t, "GetByPostAndAuthor")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Fail Lookup Existing Reaction", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(nil, errors.New("db error"))

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), constants.ReactionLike)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Database failed")
		mockReactionRepo.AssertNotCalled(t, "Create")
		mockReactionRepo.AssertNotCalled(t, "UpdateType")
		mockReactionRepo.AssertNotCalled(t, "Delete")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})

	t.Run("Fail Create", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(nil, gorm.ErrRecordNotFound)
		mockReactionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Reaction")).Return(nil, errors.New("db error"))

		res, err := u.ReactToPost(ctx, userID.String(), resID.String(), postID.String(), constants.ReactionLike)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Can not create reaction")
		mockReactionRepo.AssertNotCalled(t, "CountByPostID")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})
}

func TestReactionUseCase_FindSummaryByPostID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	resID := uuid.New()
	postID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		myReaction := &entity.Reaction{ID: uuid.New(), Type: constants.ReactionWow, PostID: postID, AuthorID: userID}
		counts := map[string]int{constants.ReactionLike: 2, constants.ReactionWow: 1}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("CountByPostID", ctx, postID.String()).Return(counts, 3, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(myReaction, nil)

		res, err := u.FindSummaryByPostID(ctx, userID.String(), resID.String(), postID.String())

		assert.NoError(t, err)
		assert.Equal(t, postID.String(), res.PostID)
		assert.Equal(t, counts, res.Counts)
		assert.Equal(t, 3, res.Total)
		assert.Equal(t, myReaction, res.MyReaction)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})

	t.Run("Success Without My Reaction", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}
		counts := map[string]int{constants.ReactionSad: 1}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("CountByPostID", ctx, postID.String()).Return(counts, 1, nil)
		mockReactionRepo.On("GetByPostAndAuthor", ctx, postID.String(), userID.String()).Return(nil, gorm.ErrRecordNotFound)

		res, err := u.FindSummaryByPostID(ctx, userID.String(), resID.String(), postID.String())

		assert.NoError(t, err)
		assert.Equal(t, 1, res.Total)
		assert.Nil(t, res.MyReaction)
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(false, nil)

		res, err := u.FindSummaryByPostID(ctx, userID.String(), resID.String(), postID.String())

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to view reactions")
		mockPostRepo.AssertNotCalled(t, "GetByID")
		mockReactionRepo.AssertNotCalled(t, "CountByPostID")
		mockUserResRepo.AssertExpectations(t)
	})

	t.Run("Fail Count", func(t *testing.T) {
		mockReactionRepo, mockPostRepo, mockUserResRepo, u := setupReactionUseCase()
		existingPost := &entity.Post{ID: postID, RestaurantID: resID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID.String()).Return(true, nil)
		mockPostRepo.On("GetByID", ctx, postID.String(), resID.String()).Return(existingPost, nil)
		mockReactionRepo.On("CountByPostID", ctx, postID.String()).Return(nil, 0, errors.New("db error"))

		res, err := u.FindSummaryByPostID(ctx, userID.String(), resID.String(), postID.String())

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Database failed")
		mockReactionRepo.AssertNotCalled(t, "GetByPostAndAuthor")
		mockUserResRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockReactionRepo.AssertExpectations(t)
	})
}
