package unit_test

import (
	"context"
	"shifty-backend/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

type MockRedisRepo struct {
	mock.Mock
}

func (m *MockRedisRepo) SaveOTP(ctx context.Context, email, otp string, duration time.Duration) error {
	args := m.Called(ctx, email, otp, duration)
	return args.Error(0)
}

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendOTP(ctx context.Context, email, otp string) error {
	args := m.Called(ctx, email, otp)
	return args.Error(0)
}
