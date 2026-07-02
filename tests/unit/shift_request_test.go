package unit_test

import (
	"context"
	"testing"
	"time"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/constants"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupShiftRequestUseCase() (*MockShiftRequestRepo, *MockUserRestaurantRepo, usecase.ShiftRequestUseCase) {
	mockShiftRequestRepo := new(MockShiftRequestRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewShiftRequestUseCase(mockShiftRequestRepo, mockUserResRepo)

	return mockShiftRequestRepo, mockUserResRepo, u
}

func TestShiftRequestUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	resID := "res-1"
	shiftID := "shift-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()

		inputRequest := &entity.ShiftRequest{
			UserID:    userID,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(8 * time.Hour),
		}
		expectedRequest := &entity.ShiftRequest{
			ID:     uuid.New(),
			UserID: userID,
			Status: constants.ShiftRequestStatusPending,
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(true, nil)
		mockShiftRequestRepo.On("Create", ctx, inputRequest).Return(expectedRequest, nil)

		res, err := u.Create(ctx, userID.String(), resID, shiftID, inputRequest)

		assert.NoError(t, err)
		assert.Equal(t, expectedRequest, res)
		mockUserResRepo.AssertExpectations(t)
		mockShiftRequestRepo.AssertExpectations(t)
	})

	t.Run("Fail Different User", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		inputRequest := &entity.ShiftRequest{UserID: uuid.New()}

		res, err := u.Create(ctx, userID.String(), resID, shiftID, inputRequest)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create shift request for another user")
		mockUserResRepo.AssertNotCalled(t, "CheckUserInRestaurant")
		mockShiftRequestRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		inputRequest := &entity.ShiftRequest{UserID: userID}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID.String(), resID).Return(false, nil)

		res, err := u.Create(ctx, userID.String(), resID, shiftID, inputRequest)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create shift request")
		mockShiftRequestRepo.AssertNotCalled(t, "Create")
	})
}

func TestShiftRequestUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	ownerID := uuid.New()
	resID := "res-1"
	shiftID := "shift-1"
	requestID := "request-1"
	statusUpdate := func() map[string]interface{} {
		return map[string]interface{}{"status": constants.ShiftRequestStatusApproved}
	}

	t.Run("Success Owner", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		existingRequest := &entity.ShiftRequest{ID: uuid.New(), UserID: userID}
		inputUpdate := map[string]interface{}{"note": "updated"}
		expectedRequest := &entity.ShiftRequest{ID: existingRequest.ID, UserID: userID}

		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(existingRequest, nil)
		mockShiftRequestRepo.On("Update", ctx, inputUpdate, shiftID, requestID).Return(expectedRequest, nil)

		res, err := u.Update(ctx, userID.String(), resID, shiftID, requestID, inputUpdate)

		assert.NoError(t, err)
		assert.Equal(t, expectedRequest, res)
		mockUserResRepo.AssertNotCalled(t, "HasManagementAuthority")
		mockShiftRequestRepo.AssertExpectations(t)
	})

	t.Run("Owner Status Ignored Without Manager Authority", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		existingRequest := &entity.ShiftRequest{ID: uuid.New(), UserID: userID}
		inputUpdate := statusUpdate()
		expectedUpdate := map[string]interface{}{}
		expectedRequest := &entity.ShiftRequest{ID: existingRequest.ID, UserID: userID}

		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(existingRequest, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(false, nil)
		mockShiftRequestRepo.On("Update", ctx, expectedUpdate, shiftID, requestID).Return(expectedRequest, nil)

		res, err := u.Update(ctx, userID.String(), resID, shiftID, requestID, inputUpdate)

		assert.NoError(t, err)
		assert.Equal(t, expectedRequest, res)
		assert.Equal(t, constants.ShiftRequestStatusApproved, inputUpdate["status"])
		mockUserResRepo.AssertExpectations(t)
		mockShiftRequestRepo.AssertExpectations(t)
	})

	t.Run("Success Manager", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		existingRequest := &entity.ShiftRequest{ID: uuid.New(), UserID: ownerID}
		inputUpdate := statusUpdate()
		expectedRequest := &entity.ShiftRequest{ID: existingRequest.ID, UserID: ownerID, Status: constants.ShiftRequestStatusApproved}

		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(existingRequest, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(true, nil)
		mockShiftRequestRepo.On("Update", ctx, inputUpdate, shiftID, requestID).Return(expectedRequest, nil)

		res, err := u.Update(ctx, userID.String(), resID, shiftID, requestID, inputUpdate)

		assert.NoError(t, err)
		assert.Equal(t, expectedRequest, res)
		mockUserResRepo.AssertExpectations(t)
		mockShiftRequestRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		existingRequest := &entity.ShiftRequest{ID: uuid.New(), UserID: ownerID}
		inputUpdate := statusUpdate()

		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(existingRequest, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(false, nil)

		res, err := u.Update(ctx, userID.String(), resID, shiftID, requestID, inputUpdate)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to update shift request")
		mockShiftRequestRepo.AssertNotCalled(t, "Update")
	})
}

func TestShiftRequestUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	ownerID := uuid.New()
	resID := "res-1"
	shiftID := "shift-1"
	requestID := "request-1"

	t.Run("Success Owner", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		existingRequest := &entity.ShiftRequest{ID: uuid.New(), UserID: userID}

		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(existingRequest, nil)
		mockShiftRequestRepo.On("Delete", ctx, shiftID, requestID).Return(nil)

		err := u.Delete(ctx, userID.String(), resID, shiftID, requestID)

		assert.NoError(t, err)
		mockUserResRepo.AssertNotCalled(t, "HasManagementAuthority")
		mockShiftRequestRepo.AssertExpectations(t)
	})

	t.Run("Success Manager", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		existingRequest := &entity.ShiftRequest{ID: uuid.New(), UserID: ownerID}

		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(existingRequest, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(true, nil)
		mockShiftRequestRepo.On("Delete", ctx, shiftID, requestID).Return(nil)

		err := u.Delete(ctx, userID.String(), resID, shiftID, requestID)

		assert.NoError(t, err)
		mockUserResRepo.AssertExpectations(t)
		mockShiftRequestRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		existingRequest := &entity.ShiftRequest{ID: uuid.New(), UserID: ownerID}

		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(existingRequest, nil)
		mockUserResRepo.On("HasManagementAuthority", ctx, userID.String(), resID).Return(false, nil)

		err := u.Delete(ctx, userID.String(), resID, shiftID, requestID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You are not allowed to delete shift request")
		mockShiftRequestRepo.AssertNotCalled(t, "Delete")
	})
}

func TestShiftRequestUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	requestID := "request-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		expectedRequest := &entity.ShiftRequest{ID: uuid.New()}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftRequestRepo.On("GetByID", ctx, shiftID, requestID).Return(expectedRequest, nil)

		res, err := u.FindByID(ctx, userID, resID, shiftID, requestID)

		assert.NoError(t, err)
		assert.Equal(t, expectedRequest, res)
	})
}

func TestShiftRequestUseCase_FindAllByShiftID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		expectedList := []*entity.ShiftRequest{
			{UserID: uuid.New()}, {UserID: uuid.New()},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftRequestRepo.On("GetAllByShiftID", ctx, shiftID).Return(expectedList, nil)

		res, err := u.FindAllByShiftID(ctx, userID, resID, shiftID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}

func TestShiftRequestUseCase_FindAllByUserID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	targetUserID := "user-2"

	t.Run("Success", func(t *testing.T) {
		mockShiftRequestRepo, mockUserResRepo, u := setupShiftRequestUseCase()
		expectedList := []*entity.ShiftRequest{
			{UserID: uuid.New()}, {UserID: uuid.New()},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftRequestRepo.On("GetAllByUserID", ctx, targetUserID).Return(expectedList, nil)

		res, err := u.FindAllByUserID(ctx, userID, resID, targetUserID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}
