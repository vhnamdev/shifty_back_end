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
				// Mock 3 bước chạy tuần tự trong Transaction
				r.On("Create", mock.Anything, mockRestaurant).Return(createdRes, nil)
				p.On("Create", mock.Anything, mock.AnythingOfType("*entity.Position")).Return(createdPos, nil)
				ur.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserRestaurant")).Return(&entity.UserRestaurant{}, nil)
			},
			expectedError: false,
		},
		{
			name:       "Fail - Invalid User ID",
			userID:     "invalid-uuid", // ID vớ vẩn, sẽ lỗi ngay vòng gửi xe
			restaurant: mockRestaurant,
			mockSetup:  func(r *MockRestaurantRepo, p *MockPositionRepo, ur *MockUserRestaurantRepo) {},
			expectedError: true,
		},
		{
			name:       "Fail - Restaurant Repo Error (Rollback)",
			userID:     userID,
			restaurant: mockRestaurant,
			mockSetup: func(r *MockRestaurantRepo, p *MockPositionRepo, ur *MockUserRestaurantRepo) {
				// Giả vờ bị sập DB ngay lúc tạo quán
				r.On("Create", mock.Anything, mockRestaurant).Return(nil, xerror.Internal("DB Error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTx := new(MockTransactor)
			mockResRepo := new(MockRestaurantRepo)
			mockPosRepo := new(MockPositionRepo)
			mockUserResRepo := new(MockUserRestaurantRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockResRepo, mockPosRepo, mockUserResRepo)
			}

			resUC := usecase.NewRestaurantUseCase(mockTx, mockResRepo, mockUserResRepo, mockPosRepo)


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