package unit_test

import (
	"context"
	"testing"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupFeedbackUseCase() (*MockFeedbackRepo, *MockUserRestaurantRepo, usecase.FeedbackUseCase) {
	mockFeedbackRepo := new(MockFeedbackRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)
	u := usecase.NewFeedbackUseCase(mockFeedbackRepo, mockUserResRepo)
	return mockFeedbackRepo, mockUserResRepo, u
}

func TestFeedbackUseCase_Create(t *testing.T) {
	ctx := context.Background()
	reviewerID := uuid.New()
	memberID := uuid.New()
	resID := uuid.New()
	t.Run("Success Manager", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		expectedFeedback := &entity.Feedback{ID: uuid.New(), Content: "Great work", RestaurantID: resID, MemberID: memberID, ReviewerID: reviewerID}
		mockUserResRepo.On("HasManagementAuthority", ctx, reviewerID.String(), resID.String()).Return(true, nil)
		mockUserResRepo.On("CheckUserInRestaurant", ctx, memberID.String(), resID.String()).Return(true, nil)
		mockFeedbackRepo.On("Create", ctx, mock.MatchedBy(func(feedback *entity.Feedback) bool {
			return feedback.Content == "Great work" && feedback.RestaurantID == resID && feedback.MemberID == memberID && feedback.ReviewerID == reviewerID
		})).Return(expectedFeedback, nil)
		res, err := u.Create(ctx, reviewerID.String(), resID.String(), memberID.String(), " Great work ")
		assert.NoError(t, err)
		assert.Equal(t, expectedFeedback, res)
		mockUserResRepo.AssertExpectations(t)
		mockFeedbackRepo.AssertExpectations(t)
	})
	t.Run("Fail Staff Forbidden", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		mockUserResRepo.On("HasManagementAuthority", ctx, reviewerID.String(), resID.String()).Return(false, nil)
		res, err := u.Create(ctx, reviewerID.String(), resID.String(), memberID.String(), "Great work")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create feedback")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockFeedbackRepo.AssertNotCalled(t, "Create")
	})
	t.Run("Fail Empty Content", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		res, err := u.Create(ctx, reviewerID.String(), resID.String(), memberID.String(), "   ")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Feedback content is required")
		mockUserResRepo.AssertNotCalled(t, "HasManagementAuthority")
		mockFeedbackRepo.AssertNotCalled(t, "Create")
	})
}

func TestFeedbackUseCase_ReadVisibility(t *testing.T) {
	ctx := context.Background()
	managerID := uuid.New()
	memberID := uuid.New()
	otherMemberID := uuid.New()
	resID := uuid.New()
	feedbackID := uuid.New()
	feedback := &entity.Feedback{ID: feedbackID, RestaurantID: resID, MemberID: memberID, ReviewerID: managerID, Content: "Solid shift"}
	t.Run("Staff Can View Own Feedback", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		mockFeedbackRepo.On("GetByID", ctx, feedbackID.String(), resID.String()).Return(feedback, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, memberID.String(), resID.String()).Return(false, nil)
		mockUserResRepo.On("CheckUserInRestaurant", ctx, memberID.String(), resID.String()).Return(true, nil)
		res, err := u.FindByID(ctx, memberID.String(), resID.String(), feedbackID.String())
		assert.NoError(t, err)
		assert.Equal(t, feedback, res)
		mockUserResRepo.AssertExpectations(t)
		mockFeedbackRepo.AssertExpectations(t)
	})
	t.Run("Staff Can Not View Other Member Feedback", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		mockFeedbackRepo.On("GetByID", ctx, feedbackID.String(), resID.String()).Return(feedback, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, otherMemberID.String(), resID.String()).Return(false, nil)
		mockUserResRepo.On("CheckUserInRestaurant", ctx, otherMemberID.String(), resID.String()).Return(true, nil)
		res, err := u.FindByID(ctx, otherMemberID.String(), resID.String(), feedbackID.String())
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to view feedback")
		mockUserResRepo.AssertExpectations(t)
		mockFeedbackRepo.AssertExpectations(t)
	})
	t.Run("Manager Can View Restaurant Feedbacks", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		feedbacks := []*entity.Feedback{feedback}
		mockUserResRepo.On("HasManagementAuthority", ctx, managerID.String(), resID.String()).Return(true, nil)
		mockFeedbackRepo.On("GetAllByRestaurantID", ctx, resID.String(), 1, 10).Return(feedbacks, int64(1), nil)
		res, total, err := u.FindAllByRestaurantID(ctx, managerID.String(), resID.String(), 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, feedbacks, res)
		mockUserResRepo.AssertExpectations(t)
		mockFeedbackRepo.AssertExpectations(t)
	})
}

func TestFeedbackUseCase_UpdateDelete(t *testing.T) {
	ctx := context.Background()
	managerID := uuid.New()
	memberID := uuid.New()
	resID := uuid.New()
	feedbackID := uuid.New()
	feedback := &entity.Feedback{ID: feedbackID, RestaurantID: resID, MemberID: memberID, ReviewerID: managerID, Content: "Old"}
	t.Run("Success Update", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		updatedFeedback := &entity.Feedback{ID: feedbackID, RestaurantID: resID, MemberID: memberID, ReviewerID: managerID, Content: "Updated"}
		updateData := map[string]interface{}{"content": "Updated"}
		mockUserResRepo.On("HasManagementAuthority", ctx, managerID.String(), resID.String()).Return(true, nil)
		mockFeedbackRepo.On("GetByID", ctx, feedbackID.String(), resID.String()).Return(feedback, nil)
		mockFeedbackRepo.On("Update", ctx, feedbackID.String(), resID.String(), updateData).Return(updatedFeedback, nil)
		res, err := u.Update(ctx, managerID.String(), resID.String(), feedbackID.String(), " Updated ")
		assert.NoError(t, err)
		assert.Equal(t, updatedFeedback, res)
		mockUserResRepo.AssertExpectations(t)
		mockFeedbackRepo.AssertExpectations(t)
	})
	t.Run("Success Soft Delete Call", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		mockUserResRepo.On("HasManagementAuthority", ctx, managerID.String(), resID.String()).Return(true, nil)
		mockFeedbackRepo.On("GetByID", ctx, feedbackID.String(), resID.String()).Return(feedback, nil)
		mockFeedbackRepo.On("Delete", ctx, feedbackID.String(), resID.String()).Return(nil)
		err := u.Delete(ctx, managerID.String(), resID.String(), feedbackID.String())
		assert.NoError(t, err)
		mockUserResRepo.AssertExpectations(t)
		mockFeedbackRepo.AssertExpectations(t)
	})
	t.Run("Fail Not Found", func(t *testing.T) {
		mockFeedbackRepo, mockUserResRepo, u := setupFeedbackUseCase()
		mockUserResRepo.On("HasManagementAuthority", ctx, managerID.String(), resID.String()).Return(true, nil)
		mockFeedbackRepo.On("GetByID", ctx, feedbackID.String(), resID.String()).Return(nil, gorm.ErrRecordNotFound)
		res, err := u.Update(ctx, managerID.String(), resID.String(), feedbackID.String(), "Updated")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Feedback is not found")
		mockFeedbackRepo.AssertNotCalled(t, "Update")
	})
}
