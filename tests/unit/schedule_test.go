package unit_test

import (
	"context"
	"errors"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/xerror"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestScheduleUseCase_Create(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	validSchedule := &entity.Schedule{StartTime: time.Now(), EndTime: time.Now().Add(8 * time.Hour)}

	tests := []struct {
		name           string
		mockAuth       func(m *MockUserRestaurantRepo)
		mockRepo       func(m *MockScheduleRepo)
		expectedError  error
		expectedResult *entity.Schedule
	}{
		{
			name: "Success",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil).Once()
			},
			mockRepo: func(m *MockScheduleRepo) {
				m.On("Create", ctx, validSchedule).Return(validSchedule, nil).Once()
			},
			expectedError:  nil,
			expectedResult: validSchedule,
		},
		{
			name: "Fail - Forbidden",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil).Once()
			},
			mockRepo:       func(m *MockScheduleRepo) {},
			expectedError:  xerror.Forbidden("You are not allowed to create schedule"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScheduleRepo := new(MockScheduleRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			tt.mockAuth(mockUserResRepo)
			tt.mockRepo(mockScheduleRepo)
			uc := usecase.NewScheduleUseCase(mockScheduleRepo, mockUserResRepo)

			result, err := uc.Create(ctx, userID, resID, validSchedule)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var expectedAppErr *xerror.AppError
				var actualAppErr *xerror.AppError
				if errors.As(tt.expectedError, &expectedAppErr) && errors.As(err, &actualAppErr) {
					assert.Equal(t, expectedAppErr.Code, actualAppErr.Code)
					assert.Equal(t, expectedAppErr.Message, actualAppErr.Message)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
			mockScheduleRepo.AssertExpectations(t)
			mockUserResRepo.AssertExpectations(t)
		})
	}
}

func TestScheduleUseCase_Update(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	scheID := uuid.New().String()
	updateData := map[string]interface{}{"number_of_members": 10}
	updatedSchedule := &entity.Schedule{NumberOfMembers: 10}

	tests := []struct {
		name           string
		mockAuth       func(m *MockUserRestaurantRepo)
		mockRepo       func(m *MockScheduleRepo)
		expectedError  error
		expectedResult *entity.Schedule
	}{
		{
			name: "Success",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil).Once()
			},
			mockRepo: func(m *MockScheduleRepo) {
				m.On("Update", ctx, resID, scheID, updateData).Return(updatedSchedule, nil).Once()
			},
			expectedError:  nil,
			expectedResult: updatedSchedule,
		},
		{
			name: "Fail - DB Update Error",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil).Once()
			},
			mockRepo: func(m *MockScheduleRepo) {
				m.On("Update", ctx, resID, scheID, updateData).Return(nil, errors.New("db error")).Once()
			},
			expectedError:  xerror.Internal("Can not update schedule"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScheduleRepo := new(MockScheduleRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			tt.mockAuth(mockUserResRepo)
			tt.mockRepo(mockScheduleRepo)
			uc := usecase.NewScheduleUseCase(mockScheduleRepo, mockUserResRepo)

			result, err := uc.Update(ctx, userID, resID, scheID, updateData)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var expectedAppErr *xerror.AppError
				var actualAppErr *xerror.AppError
				if errors.As(tt.expectedError, &expectedAppErr) && errors.As(err, &actualAppErr) {
					assert.Equal(t, expectedAppErr.Code, actualAppErr.Code)
					assert.Equal(t, expectedAppErr.Message, actualAppErr.Message)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestScheduleUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	scheID := uuid.New().String()

	tests := []struct {
		name          string
		mockAuth      func(m *MockUserRestaurantRepo)
		mockRepo      func(m *MockScheduleRepo)
		expectedError error
	}{
		{
			name: "Success",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("HasManagementAuthority", ctx, userID, resID).Return(true, nil).Once()
			},
			mockRepo: func(m *MockScheduleRepo) {
				m.On("Delete", ctx, scheID, resID).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Fail - Forbidden",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("HasManagementAuthority", ctx, userID, resID).Return(false, nil).Once()
			},
			mockRepo:      func(m *MockScheduleRepo) {},
			expectedError: xerror.Forbidden("You are not allowed to delete schedule"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScheduleRepo := new(MockScheduleRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			tt.mockAuth(mockUserResRepo)
			tt.mockRepo(mockScheduleRepo)
			uc := usecase.NewScheduleUseCase(mockScheduleRepo, mockUserResRepo)

			err := uc.Delete(ctx, userID, resID, scheID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var expectedAppErr *xerror.AppError
				var actualAppErr *xerror.AppError
				if errors.As(tt.expectedError, &expectedAppErr) && errors.As(err, &actualAppErr) {
					assert.Equal(t, expectedAppErr.Code, actualAppErr.Code)
					assert.Equal(t, expectedAppErr.Message, actualAppErr.Message)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScheduleUseCase_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	scheID := uuid.New().String()
	foundSchedule := &entity.Schedule{RestaurantID: uuid.MustParse(resID)}

	tests := []struct {
		name           string
		mockAuth       func(m *MockUserRestaurantRepo)
		mockRepo       func(m *MockScheduleRepo)
		expectedError  error
		expectedResult *entity.Schedule
	}{
		{
			name: "Success",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil).Once()
			},
			mockRepo: func(m *MockScheduleRepo) {
				m.On("FindByID", ctx, scheID, resID).Return(foundSchedule, nil).Once()
			},
			expectedError:  nil,
			expectedResult: foundSchedule,
		},
		{
			name: "Fail - Not Found (Trigger utils.IsRecordNotFoundError)",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil).Once()
			},
			mockRepo: func(m *MockScheduleRepo) {
				m.On("FindByID", ctx, scheID, resID).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedError:  xerror.NotFound("Schedule is not found"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScheduleRepo := new(MockScheduleRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			tt.mockAuth(mockUserResRepo)
			tt.mockRepo(mockScheduleRepo)
			uc := usecase.NewScheduleUseCase(mockScheduleRepo, mockUserResRepo)

			result, err := uc.FindByID(ctx, userID, resID, scheID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var expectedAppErr *xerror.AppError
				var actualAppErr *xerror.AppError
				if errors.As(tt.expectedError, &expectedAppErr) && errors.As(err, &actualAppErr) {
					assert.Equal(t, expectedAppErr.Code, actualAppErr.Code)
					assert.Equal(t, expectedAppErr.Message, actualAppErr.Message)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestScheduleUseCase_FindAllByResID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	resID := uuid.New().String()
	schedules := []*entity.Schedule{{}, {}}

	tests := []struct {
		name           string
		mockAuth       func(m *MockUserRestaurantRepo)
		mockRepo       func(m *MockScheduleRepo)
		expectedError  error
		expectedResult []*entity.Schedule
	}{
		{
			name: "Success",
			mockAuth: func(m *MockUserRestaurantRepo) {
				m.On("CheckUserInRestaurant", ctx, userID, resID).Return(true, nil).Once()
			},
			mockRepo: func(m *MockScheduleRepo) {
				m.On("FindAllByResID", ctx, resID).Return(schedules, nil).Once()
			},
			expectedError:  nil,
			expectedResult: schedules,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScheduleRepo := new(MockScheduleRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			tt.mockAuth(mockUserResRepo)
			tt.mockRepo(mockScheduleRepo)
			uc := usecase.NewScheduleUseCase(mockScheduleRepo, mockUserResRepo)

			result, err := uc.FindAllByResID(ctx, userID, resID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var expectedAppErr *xerror.AppError
				var actualAppErr *xerror.AppError
				if errors.As(tt.expectedError, &expectedAppErr) && errors.As(err, &actualAppErr) {
					assert.Equal(t, expectedAppErr.Code, actualAppErr.Code)
					assert.Equal(t, expectedAppErr.Message, actualAppErr.Message)
				} else {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}