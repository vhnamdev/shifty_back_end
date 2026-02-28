package unit_test

import (
	"context"
	"errors"
	"testing"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupPositionUseCase() (*MockPositionRepo, *MockUserRestaurantRepo, *MockTransactor, usecase.PositionUseCase) {
	mockPosRepo := new(MockPositionRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)
	mockTransactor := new(MockTransactor)

	u := usecase.NewPositionUseCase(mockPosRepo, mockUserResRepo, mockTransactor)

	return mockPosRepo, mockUserResRepo, mockTransactor, u
}

func TestPositionUseCase_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPosRepo, _, _, u := setupPositionUseCase()

		ctx := context.Background()
		inputPos := &entity.Position{Name: "Manager"}
		expectedPos := &entity.Position{ID: uuid.New(), Name: "Manager"}

		mockPosRepo.On("Create", ctx, inputPos).Return(expectedPos, nil)

		res, err := u.Create(ctx, inputPos)

		assert.NoError(t, err)
		assert.Equal(t, expectedPos, res)
		mockPosRepo.AssertExpectations(t)
	})

	t.Run("Fail Repo Error", func(t *testing.T) {
		mockPosRepo, _, _, u := setupPositionUseCase()
		ctx := context.Background()
		inputPos := &entity.Position{Name: "Manager"}

		mockPosRepo.On("Create", ctx, inputPos).Return(nil, errors.New("db error"))

		res, err := u.Create(ctx, inputPos)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Can not create position")
		mockPosRepo.AssertExpectations(t)
	})
}

func TestPositionUseCase_Update(t *testing.T) {
	ctx := context.Background()
	posID := "pos-1"
	userID := "user-1"
	resID := "res-1"
	updateData := map[string]interface{}{"name": "New Name"}

	t.Run("Success", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()
		expectedPos := &entity.Position{Name: "New Name"}

		mockUserResRepo.On("CheckAuthorityToUpdate", ctx, userID, resID).Return(true, nil)
		mockPosRepo.On("Update", ctx, posID, resID, updateData).Return(expectedPos, nil)

		res, err := u.Update(ctx, posID, userID, resID, updateData)

		assert.NoError(t, err)
		assert.Equal(t, expectedPos, res)
		mockUserResRepo.AssertExpectations(t)
		mockPosRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()

		mockUserResRepo.On("CheckAuthorityToUpdate", ctx, userID, resID).Return(false, nil)

		res, err := u.Update(ctx, posID, userID, resID, updateData)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You can not allowed to update position")

		mockPosRepo.AssertNotCalled(t, "Update", ctx, posID, resID, updateData)
	})
}

func TestPositionUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	posID := "pos-1"
	userID := "user-1"
	resID := "res-1"

	t.Run("Success", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()

		mockUserResRepo.On("CheckAuthorityToDelete", ctx, userID, resID).Return(true, nil)
		mockUserResRepo.On("SetPositionNull", ctx, posID, resID).Return(nil)
		mockPosRepo.On("Delete", ctx, posID, resID).Return(nil)

		err := u.Delete(ctx, userID, resID, posID)

		assert.NoError(t, err)
		mockUserResRepo.AssertExpectations(t)
		mockPosRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()

		mockUserResRepo.On("CheckAuthorityToDelete", ctx, userID, resID).Return(false, nil)

		err := u.Delete(ctx, userID, resID, posID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You can not allowed to delete position")

		mockUserResRepo.AssertNotCalled(t, "SetPositionNull")
		mockPosRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Fail Transaction Error", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()

		mockUserResRepo.On("CheckAuthorityToDelete", ctx, userID, resID).Return(true, nil)
		mockUserResRepo.On("SetPositionNull", ctx, posID, resID).Return(errors.New("db error"))

		err := u.Delete(ctx, userID, resID, posID)

		assert.Error(t, err)
		mockPosRepo.AssertNotCalled(t, "Delete")
	})
}

func TestPositionUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	posID := "pos-1"
	userID := "user-1"
	resID := "res-1"

	t.Run("Success", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()
		expectedPos := &entity.Position{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPosRepo.On("FindByID", ctx, posID, resID).Return(expectedPos, nil)

		res, err := u.FindByID(ctx, posID, userID, resID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPos, res)
	})

	t.Run("Fail Not Member", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(false, nil)

		res, err := u.FindByID(ctx, posID, userID, resID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed of this restaurant")

		mockPosRepo.AssertNotCalled(t, "FindByID")
	})
}

func TestPositionUseCase_GetAllByRestaurantID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"

	t.Run("Success", func(t *testing.T) {
		mockPosRepo, mockUserResRepo, _, u := setupPositionUseCase()
		expectedList := []*entity.Position{
			{Name: "A"}, {Name: "B"},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockPosRepo.On("GetAllByRestaurantID", ctx, resID).Return(expectedList, nil)

		res, err := u.GetAllByRestaurantID(ctx, resID, userID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}
