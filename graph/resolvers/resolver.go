//go:generate go run github.com/99designs/gqlgen generate
package resolvers

import (
	"shifty-backend/graph"
	"shifty-backend/internal/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	UserUseCase             usecase.UserUseCase
	UserRestaurantUseCase   usecase.UserRestaurantUseCase
	RestaurantUseCase       usecase.RestaurantUseCase
	PositionUseCase         usecase.PositionUseCase
	ScheduleUseCase         usecase.ScheduleUseCase
	ShiftUseCase            usecase.ShiftUseCase
	ShiftRequirementUseCase usecase.ShiftRequirementUseCase
	ShiftRuleUseCase        usecase.ShiftRuleUseCase
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
