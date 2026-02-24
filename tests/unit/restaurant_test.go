package unit_test

import (
	"context"
	"testing"

	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/xerror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRestaurantUseCase_Create(t *testing.T) {
	userID := uuid.New().String()
	resID := uuid.New()
	posID := uuid.New()

	mockRestaurant := &entity.Restaurant{Name: "Quán Nướng F&B"}
	createdRes := &entity.Restaurant{ID: resID, Name: "Quán Nướng F&B"}
	createdPos := &entity.Position{ID: posID, Name: "OWNER"}

	tests := []struct {
		name          string
		userID        string
		restaurant    *entity.Restaurant
		mockSetup     func(*MockRestaurantRepo, *MockPositionRepo, *MockUserRestaurantRepo)
		expectedError bool
	}{
		{
			name:       "Success - Create Restaurant with Owner",
			userID:     userID,
			restaurant: mockRestaurant,
			mockSetup: func(r *MockRestaurantRepo, p *MockPositionRepo, ur *MockUserRestaurantRepo) {
				r.On("Create", mock.Anything, mockRestaurant).Return(createdRes, nil)
				p.On("Create", mock.Anything, mock.AnythingOfType("*entity.Position")).Return(createdPos, nil)
				ur.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserRestaurant")).Return(&entity.UserRestaurant{}, nil)
			},
			expectedError: false,
		},
		{
			name:          "Fail - Invalid User ID",
			userID:        "invalid-uuid",
			restaurant:    mockRestaurant,
			mockSetup:     func(r *MockRestaurantRepo, p *MockPositionRepo, ur *MockUserRestaurantRepo) {},
			expectedError: true,
		},
		{
			name:       "Fail - Restaurant Repo Error (Rollback)",
			userID:     userID,
			restaurant: mockRestaurant,
			mockSetup: func(r *MockRestaurantRepo, p *MockPositionRepo, ur *MockUserRestaurantRepo) {
				r.On("Create", mock.Anything, mockRestaurant).Return(nil, xerror.Internal("DB Error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTx := new(MockTransactor)
			mockResRepo := new(MockRestaurantRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)
			mockPosRepo := new(MockPositionRepo)
			mockRedisRepo := new(MockRedisRepo)
			mockUserRepo := new(MockUserRepo)
			mockMailer := new(MockMailer)
			mockUploader := new(MockUploader)

			if tt.mockSetup != nil {
				tt.mockSetup(mockResRepo, mockPosRepo, mockUserResRepo)
			}

			resUC := usecase.NewRestaurantUseCase(
				mockTx,
				mockResRepo,
				mockUserResRepo,
				mockPosRepo,
				mockRedisRepo,
				mockUserRepo,
				mockMailer,
				mockUploader,
			)

			result, err := resUC.Create(context.Background(), tt.userID, tt.restaurant)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockResRepo.AssertExpectations(t)
			mockPosRepo.AssertExpectations(t)
			mockUserResRepo.AssertExpectations(t)
		})
	}
}
