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

func setupShiftRequirementUseCase() (*MockShiftRequirementRepo, *MockUserRestaurantRepo, usecase.ShiftRequirementUseCase) {
	mockShiftReqRepo := new(MockShiftRequirementRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewShiftRequirementUseCase(mockShiftReqRepo, mockUserResRepo)

	return mockShiftReqRepo, mockUserResRepo, u
}

func TestShiftRequirementUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()

		inputReq := &entity.ShiftRequirement{Quantity: 2, StartTime: time.Now()}
		expectedReq := &entity.ShiftRequirement{ID: uuid.New(), Quantity: 2}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftReqRepo.On("Create", ctx, inputReq).Return(expectedReq, nil)

		res, err := u.Create(ctx, userID, resID, shiftID, inputReq)

		assert.NoError(t, err)
		assert.Equal(t, expectedReq, res)
		mockUserResRepo.AssertExpectations(t)
		mockShiftReqRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()
		inputReq := &entity.ShiftRequirement{Quantity: 2}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Create(ctx, userID, resID, shiftID, inputReq)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create shift")
		mockShiftReqRepo.AssertNotCalled(t, "Create")
	})
}

func TestShiftRequirementUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	reqID := "req-1"
	updateData := map[string]interface{}{"quantity": 3}

	t.Run("Success", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()
		expectedReq := &entity.ShiftRequirement{Quantity: 3}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftReqRepo.On("Update", ctx, updateData, shiftID, reqID).Return(expectedReq, nil)

		res, err := u.Update(ctx, userID, resID, shiftID, reqID, updateData)

		assert.NoError(t, err)
		assert.Equal(t, expectedReq, res)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Update(ctx, userID, resID, shiftID, reqID, updateData)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create shift")
		mockShiftReqRepo.AssertNotCalled(t, "Update")
	})
}

func TestShiftRequirementUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	reqID := "req-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftReqRepo.On("Delete", ctx, shiftID, reqID).Return(nil)

		err := u.Delete(ctx, userID, resID, shiftID, reqID)

		assert.NoError(t, err)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		err := u.Delete(ctx, userID, resID, shiftID, reqID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You are not allowed to create shift")
		mockShiftReqRepo.AssertNotCalled(t, "Delete")
	})
}

func TestShiftRequirementUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"
	reqID := "req-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()
		expectedReq := &entity.ShiftRequirement{ID: uuid.New()}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftReqRepo.On("GetByID", ctx, shiftID, reqID).Return(expectedReq, nil)

		res, err := u.FindByID(ctx, userID, resID, shiftID, reqID)

		assert.NoError(t, err)
		assert.Equal(t, expectedReq, res)
	})
}

func TestShiftRequirementUseCase_FindAllByShiftID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	shiftID := "shift-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftReqRepo, mockUserResRepo, u := setupShiftRequirementUseCase()
		expectedList := []*entity.ShiftRequirement{
			{Quantity: 1}, {Quantity: 2},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftReqRepo.On("GetAllByShiftID", ctx, shiftID).Return(expectedList, nil)

		res, err := u.FindAllByShiftID(ctx, userID, resID, shiftID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}
