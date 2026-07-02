package usecase

import (
	"context"
	"encoding/json"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"

	"gorm.io/datatypes"
)

type ShiftRuleUseCase interface {
	Create(ctx context.Context, userID, resID string, shiftRule *entity.ShiftRule) (*entity.ShiftRule, error)
	Update(ctx context.Context, userID, resID string, updateData map[string]interface{}, shiftRuleID, restaurantID string) (*entity.ShiftRule, error)
	Delete(ctx context.Context, userID, resID string, shiftRuleID, restaurantID string) error
	FindByID(ctx context.Context, userID, resID string, shiftRuleID, restaurantID string) (*entity.ShiftRule, error)
	FindAllByRestaurantID(ctx context.Context, userID, resID string, restaurantID string) ([]*entity.ShiftRule, error)
}

type shiftRuleUseCase struct {
	shiftRuleRepo      repository.ShiftRuleRepository
	userRestaurantRepo repository.UserRestaurantRepository
}

func NewShiftRuleUseCase(shiftRuleRepo repository.ShiftRuleRepository, userRestaurantRepo repository.UserRestaurantRepository) ShiftRuleUseCase {
	return &shiftRuleUseCase{
		shiftRuleRepo:      shiftRuleRepo,
		userRestaurantRepo: userRestaurantRepo,
	}
}

// formatRuleConfig validates and rebuilds a clean, safe JSON config based on the rule type
func formatRuleConfig(ruleType string, rawConfig datatypes.JSON) (datatypes.JSON, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(rawConfig, &rawData); err != nil {
		return nil, xerror.BadRequest("Invalid config data format")
	}

	safeMap := make(map[string]interface{})

	switch ruleType {
	case constants.RuleTypeMaxHoursPerDay, constants.RuleTypeMaxHoursPerWeek:
		val, ok := rawData["hours"]
		if !ok {
			return nil, xerror.BadRequest("Missing 'hours' parameter for this rule")
		}
		safeMap["max_hours"] = val

	case constants.RuleTypeMinRestTime:
		val, ok := rawData["hours"]
		if !ok {
			return nil, xerror.BadRequest("Missing 'hours' parameter for this rule")
		}
		safeMap["min_hours"] = val

	case constants.RuleTypeMustWorkWith, constants.RuleTypeBanWorkWith:
		val, ok := rawData["user_ids"]
		if !ok {
			return nil, xerror.BadRequest("Missing 'user_ids' parameter for this rule")
		}
		safeMap["user_ids"] = val

	case constants.RuleTypeQualification:
		val, ok := rawData["position_id"]
		if !ok {
			return nil, xerror.BadRequest("Missing 'position_id' parameter for this rule")
		}
		safeMap["position_id"] = val

	default:
		return nil, xerror.BadRequest("Unsupported rule type")
	}

	safeJSONBytes, err := json.Marshal(safeMap)
	if err != nil {
		return nil, xerror.Internal("Internal error while formatting config")
	}

	return datatypes.JSON(safeJSONBytes), nil
}

func (u *shiftRuleUseCase) Create(ctx context.Context, userID, resID string, shiftRule *entity.ShiftRule) (*entity.ShiftRule, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to create shift rule")
	}

	formattedConfig, err := formatRuleConfig(shiftRule.Type, shiftRule.Config)
	if err != nil {
		return nil, err
	}
	shiftRule.Config = formattedConfig

	newShiftRule, err := u.shiftRuleRepo.Create(ctx, shiftRule)

	if err != nil {
		return nil, xerror.Internal("Can not create shift rule")
	}

	return newShiftRule, nil
}

func (u *shiftRuleUseCase) Update(ctx context.Context, userID, resID string, updateData map[string]interface{}, shiftRuleID, restaurantID string) (*entity.ShiftRule, error) {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to update shift rule")
	}

	_, hasType := updateData["type"]
	rawConfig, hasConfig := updateData["config"]

	if hasType && !hasConfig {
		return nil, xerror.BadRequest("You must provide a new 'config' when changing the rule 'type'")
	}

	if hasConfig {
		var ruleType string

		if rt, ok := updateData["type"].(string); ok {
			ruleType = rt
		} else {

			existingRule, err := u.shiftRuleRepo.GetByID(ctx, shiftRuleID, restaurantID)
			if err != nil {
				return nil, xerror.NotFound("Shift rule not found for config validation")
			}
			ruleType = existingRule.Type
		}

		formattedConfig, err := formatRuleConfig(ruleType, rawConfig.(datatypes.JSON))
		if err != nil {
			return nil, err
		}
		updateData["config"] = formattedConfig
	}

	updatedShiftRule, err := u.shiftRuleRepo.Update(ctx, updateData, shiftRuleID, restaurantID)

	if err != nil {
		return nil, xerror.Internal("Can not update shift rule")
	}

	return updatedShiftRule, nil
}

func (u *shiftRuleUseCase) Delete(ctx context.Context, userID, resID string, shiftRuleID, restaurantID string) error {
	isAuthority, err := u.userRestaurantRepo.HasManagementAuthority(ctx, userID, resID)

	if err != nil {
		return xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return xerror.Forbidden("You are not allowed to delete shift rule")
	}

	return u.shiftRuleRepo.Delete(ctx, shiftRuleID, restaurantID)
}

func (u *shiftRuleUseCase) FindByID(ctx context.Context, userID, resID string, shiftRuleID, restaurantID string) (*entity.ShiftRule, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view shift rule")
	}

	shiftRule, err := u.shiftRuleRepo.GetByID(ctx, shiftRuleID, restaurantID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift rule is not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftRule, nil
}

func (u *shiftRuleUseCase) FindAllByRestaurantID(ctx context.Context, userID, resID string, restaurantID string) ([]*entity.ShiftRule, error) {
	isAuthority, err := u.userRestaurantRepo.CheckUserInRestaurant(ctx, userID, resID)

	if err != nil {
		return nil, xerror.Internal("Can not check authority")
	}

	if !isAuthority {
		return nil, xerror.Forbidden("You are not allowed to view shift rules")
	}

	shiftRules, err := u.shiftRuleRepo.GetAllByRestaurantID(ctx, restaurantID)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("Shift rules are not found")
		}
		return nil, xerror.Internal("Database failed")
	}

	return shiftRules, nil
}
