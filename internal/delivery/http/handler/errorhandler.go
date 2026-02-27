package handler

import (
	"errors"
	"shifty-backend/pkg/xerror"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

// GlobalErrorHandler centralized error handling for all services and controllers.
// This prevents the need for repetitive error checks in every function.
func GlobalErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *xerror.AppError
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}
	if code >= 500 {
		if hub := sentryfiber.GetHubFromContext(ctx); hub != nil {
			hub.CaptureException(err)
		} else {
			sentry.CaptureException(err)
		}
	}
	return ctx.Status(code).JSON(fiber.Map{
		"status":  "Error",
		"code":    code,
		"message": message,
	})
}
