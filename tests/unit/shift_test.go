package unit_test

import (
	"context"
	"testing"
	"time"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupShiftUseCase() (*MockShiftRepo, *MockUserRestaurantRepo, usecase.ShiftUseCase) {
	mockShiftRepo := new(MockShiftRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewShiftUseCase(mockShiftRepo, mockUserResRepo)

	return mockShiftRepo, mockUserResRepo, u
}

func TestShiftUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()

		inputShift := &entity.Shift{Type: "Morning", StartTime: time.Now()}
		expectedShift := &entity.Shift{ID: uuid.New(), Type: "Morning"}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftRepo.On("Create", ctx, inputShift).Return(expectedShift, nil)

		res, err := u.Create(ctx, userID, resID, inputShift)

		assert.NoError(t, err)
		assert.Equal(t, expectedShift, res)
		mockUserResRepo.AssertExpectations(t)
		mockShiftRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()
		inputShift := &entity.Shift{Type: "Morning"}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Create(ctx, userID, resID, inputShift)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create shift")
		mockShiftRepo.AssertNotCalled(t, "Create")
	})
}

func TestShiftUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	scheID := "sche-1"
	shiftID := "shift-1"
	updateData := map[string]interface{}{"type": "Evening"}

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()
		expectedShift := &entity.Shift{Type: "Evening"}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftRepo.On("Update", ctx, shiftID, scheID, updateData).Return(expectedShift, nil)

		res, err := u.Update(ctx, userID, resID, scheID, shiftID, updateData)

		assert.NoError(t, err)
		assert.Equal(t, expectedShift, res)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Update(ctx, userID, resID, scheID, shiftID, updateData)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to update shift")
		mockShiftRepo.AssertNotCalled(t, "Update")
	})
}

func TestShiftUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	scheID := "sche-1"
	shiftID := "shift-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftRepo.On("Delete", ctx, scheID, shiftID).Return(nil)

		err := u.Delete(ctx, userID, resID, scheID, shiftID)

		assert.NoError(t, err)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		err := u.Delete(ctx, userID, resID, scheID, shiftID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You are not allowed to delete shift")
		mockShiftRepo.AssertNotCalled(t, "Delete")
	})
}

func TestShiftUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	scheID := "sche-1"
	shiftID := "shift-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()
		expectedShift := &entity.Shift{ID: uuid.New()}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftRepo.On("FindByID", ctx, scheID, shiftID).Return(expectedShift, nil)

		res, err := u.FindByID(ctx, userID, resID, scheID, shiftID)

		assert.NoError(t, err)
		assert.Equal(t, expectedShift, res)
	})
}

func TestShiftUseCase_FindAllByScheduleID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	scheID := "sche-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo, mockUserResRepo, u := setupShiftUseCase()
		expectedList := []*entity.Shift{
			{Type: "Morning"}, {Type: "Evening"},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftRepo.On("FindAllByScheduleID", ctx, scheID).Return(expectedList, nil)

		res, err := u.FindAllByScheduleID(ctx, userID, resID, scheID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}
