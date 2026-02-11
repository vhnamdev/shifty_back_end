package handler

import (
	"shifty-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUC usecase.UserUseCase
}

func NewUserHandler(userUC usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
	}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error{
	return nil
}