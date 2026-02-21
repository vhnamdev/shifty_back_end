package repository

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

const txKey = contextKey("txKey")

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactorImpl struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) Transactor {
	return &transactorImpl{
		db: db,
	}
}

func (t *transactorImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		
		return fn(txCtx)
	})
}

func Extract(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}

	return defaultDB
}
