package handler

import (
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/mailer"
	"shifty-backend/pkg/uploader"
	"shifty-backend/pkg/xerror"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUC   usecase.AuthUseCase
	uploader uploader.ImageUploader
	mailer   mailer.EmailSender
}

func NewAuthHandler(authUC usecase.AuthUseCase, uploader uploader.ImageUploader, mailer mailer.EmailSender) *AuthHandler {
	return &AuthHandler{
		authUC:   authUC,
		uploader: uploader,
		mailer:   mailer,
	}
}

// Func Register account
func (h *AuthHandler) RegisterLocal(c *fiber.Ctx) error {

	//Parse request body to struct RegisterRequest
	req := new(dto.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return xerror.BadRequest("Invalid request body")
	}
	ctx := c.UserContext()

	file, err := c.FormFile("avatar")
	if err != nil {
		return xerror.BadRequest("Please choose image to upload")
	}
	imageURL, errUpload := h.uploader.UploadImage(ctx, file, "avatars")
	if errUpload != nil {
		return xerror.Internal("Fail to save avatar into Cloudinary")
	}
	newUser := &entity.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Avatar:   imageURL,
		Role:     constants.User,
	}
	err = h.authUC.RegisterLocal(ctx, newUser)

	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"message": "OK"})
}

// LoginLocal handles authentication with Email & Password
func (h *AuthHandler) LoginLocal(c *fiber.Ctx) error {

	// Get Login Request
	req := new(dto.LoginRequest)

	// Parse & Validate request body
	if err := c.BodyParser(req); err != nil {
		return xerror.BadRequest("Invalid request body")
	}

	// Create context
	ctx := c.UserContext()

	// Extract client info (for security logging)
	ua := c.Get("User-Agent")
	// Get user's device IP
	ip := c.IP()

	// Call function Login Local in usecase
	accessToken, refreshToken, user, err := h.authUC.LoginLocal(ctx, req.Email, req.Password, ua, ip)

	if err != nil {
		return err
	}

	// Set Refresh Token to Secure Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})
	return c.Status(200).JSON(fiber.Map{
		"access_token": accessToken,
		"user":         user,
	})
}

// Login with Google
func (h *AuthHandler) LoginGoogle(c *fiber.Ctx) error {

	// Get login google request
	req := new(dto.GoogleLogin)

	ctx := c.UserContext()

	ua := c.Get("User-Agent")
	ip := c.IP()

	if err := c.BodyParser(req); err != nil {
		return xerror.BadRequest("Invalid request")
	}

	accessToken, refreshToken, user, err := h.authUC.LoginGoogle(ctx, req.Code, ua, ip)

	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Status(200).JSON(fiber.Map{
		"access_token": accessToken,
		"user":         user,
	})
}

// Func send otp to user's email
func (h *AuthHandler) SendOTP(c *fiber.Ctx) error {

	// Get sendotp request
	req := new(dto.SendOTP)
	if err := c.BodyParser(req); err != nil {
		return xerror.BadRequest("Invalid Email")
	}

	ctx := c.UserContext()

	// Call sendOTP usecase
	err := h.authUC.SendOTP(ctx, req.Email, req.Purpose)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": "OK"})
}

// Reset password
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {

	req := new(dto.ResetPassword)

	if err := c.BodyParser(req); err != nil {
		return xerror.BadRequest("Invalid password")
	}

	ctx := c.UserContext()

	// Call func reset password
	err := h.authUC.ResetPassword(ctx, req.Email, req.Password)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": "OK"})
}
