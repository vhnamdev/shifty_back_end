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

func setupShiftAssignmentUseCase() (*MockShiftAssignmentRepo, *MockUserRestaurantRepo, usecase.ShiftAssignmentUseCase) {
	mockShiftAssignmentRepo := new(MockShiftAssignmentRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewShiftAssignmentUseCase(mockShiftAssignmentRepo, mockUserResRepo)

	return mockShiftAssignmentRepo, mockUserResRepo, u
}

func TestShiftAssignmentUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	staffID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()

		inputAssignment := &entity.ShiftAssignment{UserID: staffID, CheckInTime: timePtr(time.Now())}
		expectedAssignment := &entity.ShiftAssignment{ID: uuid.New(), UserID: staffID}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftAssignmentRepo.On("IsExistByUserAndShift", ctx, shiftID, staffID.String()).Return(false, nil)
		mockShiftAssignmentRepo.On("Create", ctx, inputAssignment).Return(expectedAssignment, nil)

		res, err := u.Create(ctx, userID, resID, shiftID, inputAssignment)

		assert.NoError(t, err)
		assert.Equal(t, expectedAssignment, res)
		mockUserResRepo.AssertExpectations(t)
		mockShiftAssignmentRepo.AssertExpectations(t)
	})

	t.Run("Fail Duplicate Staff", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()
		inputAssignment := &entity.ShiftAssignment{UserID: staffID}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftAssignmentRepo.On("IsExistByUserAndShift", ctx, shiftID, staffID.String()).Return(true, nil)

		res, err := u.Create(ctx, userID, resID, shiftID, inputAssignment)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Staff already has assignment in this shift")
		mockShiftAssignmentRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()
		inputAssignment := &entity.ShiftAssignment{UserID: staffID}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Create(ctx, userID, resID, shiftID, inputAssignment)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create shift assignment")
		mockShiftAssignmentRepo.AssertNotCalled(t, "IsExistByUserAndShift")
		mockShiftAssignmentRepo.AssertNotCalled(t, "Create")
	})
}

func TestShiftAssignmentUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	assignmentID := "assignment-1"
	updateData := map[string]interface{}{"note": "Updated note"}

	t.Run("Success", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()
		expectedAssignment := &entity.ShiftAssignment{Note: stringPtr("Updated note")}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftAssignmentRepo.On("Update", ctx, updateData, shiftID, assignmentID).Return(expectedAssignment, nil)

		res, err := u.Update(ctx, userID, resID, shiftID, assignmentID, updateData)

		assert.NoError(t, err)
		assert.Equal(t, expectedAssignment, res)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Update(ctx, userID, resID, shiftID, assignmentID, updateData)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to update shift assignment")
		mockShiftAssignmentRepo.AssertNotCalled(t, "Update")
	})
}

func TestShiftAssignmentUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	assignmentID := "assignment-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftAssignmentRepo.On("Delete", ctx, shiftID, assignmentID).Return(nil)

		err := u.Delete(ctx, userID, resID, shiftID, assignmentID)

		assert.NoError(t, err)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		err := u.Delete(ctx, userID, resID, shiftID, assignmentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You are not allowed to delete shift assignment")
		mockShiftAssignmentRepo.AssertNotCalled(t, "Delete")
	})
}

func TestShiftAssignmentUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	assignmentID := "assignment-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()
		expectedAssignment := &entity.ShiftAssignment{ID: uuid.New()}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftAssignmentRepo.On("GetByID", ctx, shiftID, assignmentID).Return(expectedAssignment, nil)

		res, err := u.FindByID(ctx, userID, resID, shiftID, assignmentID)

		assert.NoError(t, err)
		assert.Equal(t, expectedAssignment, res)
	})
}

func TestShiftAssignmentUseCase_FindAllByShiftID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftAssignmentRepo, mockUserResRepo, u := setupShiftAssignmentUseCase()
		expectedList := []*entity.ShiftAssignment{
			{UserID: uuid.New()}, {UserID: uuid.New()},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftAssignmentRepo.On("GetAllByShiftID", ctx, shiftID).Return(expectedList, nil)

		res, err := u.FindAllByShiftID(ctx, userID, resID, shiftID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func stringPtr(s string) *string {
	return &s
}
