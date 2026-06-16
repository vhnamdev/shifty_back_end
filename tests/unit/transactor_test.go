package unit_test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK TRANSACTOR =====
type MockTransactor struct{ mock.Mock }

func (m *MockTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
