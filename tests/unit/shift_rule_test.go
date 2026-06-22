package unit_test

import (
	"context"
	"testing"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func setupShiftRuleUseCase() (*MockShiftRuleRepo, *MockUserRestaurantRepo, usecase.ShiftRuleUseCase) {
	mockShiftRuleRepo := new(MockShiftRuleRepo)
	mockUserResRepo := new(MockUserRestaurantRepo)

	u := usecase.NewShiftRuleUseCase(mockShiftRuleRepo, mockUserResRepo)

	return mockShiftRuleRepo, mockUserResRepo, u
}

func TestShiftRuleUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()

		inputShiftRule := &entity.ShiftRule{
			Type:   entity.RuleTypeMaxHoursPerDay,
			Name:   "Max hours per day",
			Config: datatypes.JSON(`{"hours":8,"ignored":"value"}`),
		}
		expectedShiftRule := &entity.ShiftRule{
			ID:     uuid.New(),
			Type:   entity.RuleTypeMaxHoursPerDay,
			Name:   "Max hours per day",
			Config: datatypes.JSON(`{"max_hours":8}`),
		}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftRuleRepo.On("Create", ctx, inputShiftRule).Return(expectedShiftRule, nil)

		res, err := u.Create(ctx, userID, resID, inputShiftRule)

		assert.NoError(t, err)
		assert.Equal(t, expectedShiftRule, res)
		assert.JSONEq(t, `{"max_hours":8}`, string(inputShiftRule.Config))
		mockUserResRepo.AssertExpectations(t)
		mockShiftRuleRepo.AssertExpectations(t)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()
		inputShiftRule := &entity.ShiftRule{
			Type:   entity.RuleTypeMaxHoursPerDay,
			Name:   "Max hours per day",
			Config: datatypes.JSON(`{"hours":8}`),
		}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Create(ctx, userID, resID, inputShiftRule)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to create shift rule")
		mockShiftRuleRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Fail Invalid Config", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()
		inputShiftRule := &entity.ShiftRule{
			Type:   entity.RuleTypeMaxHoursPerDay,
			Name:   "Max hours per day",
			Config: datatypes.JSON(`{"max_hours":8}`),
		}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)

		res, err := u.Create(ctx, userID, resID, inputShiftRule)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Missing 'hours' parameter for this rule")
		mockShiftRuleRepo.AssertNotCalled(t, "Create")
	})
}

func TestShiftRuleUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	restaurantID := "restaurant-1"
	shiftRuleID := "shift-rule-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()
		updateData := map[string]interface{}{
			"type":   entity.RuleTypeMinRestTime,
			"config": datatypes.JSON(`{"hours":12}`),
		}
		expectedShiftRule := &entity.ShiftRule{Type: entity.RuleTypeMinRestTime}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftRuleRepo.On("Update", ctx, updateData, shiftRuleID, restaurantID).Return(expectedShiftRule, nil)

		res, err := u.Update(ctx, userID, resID, updateData, shiftRuleID, restaurantID)

		assert.NoError(t, err)
		assert.Equal(t, expectedShiftRule, res)
		assert.JSONEq(t, `{"min_hours":12}`, string(updateData["config"].(datatypes.JSON)))
	})

	t.Run("Success With Existing Rule Type", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()
		updateData := map[string]interface{}{
			"config": datatypes.JSON(`{"user_ids":["user-2","user-3"]}`),
		}
		existingShiftRule := &entity.ShiftRule{Type: entity.RuleTypeMustWorkWith}
		expectedShiftRule := &entity.ShiftRule{Type: entity.RuleTypeMustWorkWith}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftRuleRepo.On("GetByID", ctx, shiftRuleID, restaurantID).Return(existingShiftRule, nil)
		mockShiftRuleRepo.On("Update", ctx, updateData, shiftRuleID, restaurantID).Return(expectedShiftRule, nil)

		res, err := u.Update(ctx, userID, resID, updateData, shiftRuleID, restaurantID)

		assert.NoError(t, err)
		assert.Equal(t, expectedShiftRule, res)
		assert.JSONEq(t, `{"user_ids":["user-2","user-3"]}`, string(updateData["config"].(datatypes.JSON)))
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()
		updateData := map[string]interface{}{"name": "Updated rule"}

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		res, err := u.Update(ctx, userID, resID, updateData, shiftRuleID, restaurantID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "You are not allowed to update shift rule")
		mockShiftRuleRepo.AssertNotCalled(t, "Update")
	})
}

func TestShiftRuleUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	restaurantID := "restaurant-1"
	shiftRuleID := "shift-rule-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil)
		mockShiftRuleRepo.On("Delete", ctx, shiftRuleID, restaurantID).Return(nil)

		err := u.Delete(ctx, userID, resID, shiftRuleID, restaurantID)

		assert.NoError(t, err)
	})

	t.Run("Fail Forbidden", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()

		mockUserResRepo.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil)

		err := u.Delete(ctx, userID, resID, shiftRuleID, restaurantID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You are not allowed to delete shift rule")
		mockShiftRuleRepo.AssertNotCalled(t, "Delete")
	})
}

func TestShiftRuleUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	restaurantID := "restaurant-1"
	shiftRuleID := "shift-rule-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()
		expectedShiftRule := &entity.ShiftRule{ID: uuid.New()}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftRuleRepo.On("GetByID", ctx, shiftRuleID, restaurantID).Return(expectedShiftRule, nil)

		res, err := u.FindByID(ctx, userID, resID, shiftRuleID, restaurantID)

		assert.NoError(t, err)
		assert.Equal(t, expectedShiftRule, res)
	})
}

func TestShiftRuleUseCase_FindAllByRestaurantID(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	resID := "res-1"
	restaurantID := "restaurant-1"

	t.Run("Success", func(t *testing.T) {
		mockShiftRuleRepo, mockUserResRepo, u := setupShiftRuleUseCase()
		expectedList := []*entity.ShiftRule{
			{Type: entity.RuleTypeMaxHoursPerDay}, {Type: entity.RuleTypeMinRestTime},
		}

		mockUserResRepo.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil)
		mockShiftRuleRepo.On("GetAllByRestaurantID", ctx, restaurantID).Return(expectedList, nil)

		res, err := u.FindAllByRestaurantID(ctx, userID, resID, restaurantID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}
