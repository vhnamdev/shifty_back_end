package handler

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
	}
}

func (h *AuthHandler) LoginLocal(ctx context.Context, email string, password string) (string, *entity.User, error) {
	return h.authUC.LoginLocal(ctx, email, password, "", "")
}
