package handler

import (
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/uploader"
	"shifty-backend/pkg/xerror"

	"github.com/gofiber/fiber/v2"
)

type RestaurantHandler struct {
	restaurantUC usecase.RestaurantUseCase
	uploader     uploader.CloudinaryService
}

func NewRestaurantHandler(restaurantUC usecase.RestaurantUseCase, uploader uploader.CloudinaryService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantUC: restaurantUC,
		uploader:     uploader,
	}
}

func (h *RestaurantHandler) UpdateImage(c *fiber.Ctx) error {
	resID := c.Params("id")
	if resID == "" {
		return xerror.BadRequest("Restaurant ID is required")
	}

	ctx := c.UserContext()

	userID, ok := ctx.Value("user_id").(string)

	if !ok || userID == "" {
		return xerror.BadRequest("You are not logged in")
	}
	file, err := c.FormFile("avatar")
	if err != nil {
		return xerror.BadRequest("Please choose image to upload")
	}
	imageURL, errUpload := h.uploader.UploadImage(ctx, file, "restaurants")
	if errUpload != nil {
		return xerror.Internal("Fail to save avatar into Cloudinary")
	}

	updatedRestaurant, err := h.restaurantUC.UpdateImage(ctx, userID, resID, imageURL)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Update avatar successful",
		"data":    updatedRestaurant,
	})
}
