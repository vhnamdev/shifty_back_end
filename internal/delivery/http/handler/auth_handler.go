package handler

import (
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/mailer"
	"shifty-backend/pkg/uploader"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUC   usecase.AuthUseCase
	uploader *uploader.CloudinaryService
	mailer   *mailer.EmailService
}

func NewAuthHandler(authUC usecase.AuthUseCase, uploader *uploader.CloudinaryService, mailer *mailer.EmailService) *AuthHandler {
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
	ctx := c.Context()

	avatarURL := ""

	file, err := c.FormFile("avatar")

	if err == nil {
		url, errUpload := h.uploader.UploadImage(ctx, file)
		if errUpload != nil {
			return xerror.BadRequest("Fail to save avatar into Cloudinary")
		}
		avatarURL = url
	}
	newUser := &entity.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Avatar:   avatarURL,
		Role:     "User",
	}
	err = h.authUC.RegisterLocal(ctx, newUser)

	if err != nil {
		return err
	}
	return c.Status(200).JSON("OK")
}

// Login with Emal and Password
func (h *AuthHandler) LoginLocal(c *fiber.Ctx) error {

	// Get Login Request
	req := new(dto.LoginRequest)

	//
	if err := c.BodyParser(req); err != nil {
		return xerror.BadRequest("Invalid request body")
	}
	ctx := c.UserContext()

	ua := c.Get("User-Agent")
	ip := c.IP()

	at, rt, user, err := h.authUC.LoginLocal(ctx, req.Email, req.Password, ip, ua)

	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    rt,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})
	return c.Status(200).JSON(fiber.Map{
		"access_token": at,
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

	// Find user by email
	user, err := h.authUC.FindUserByEmail(ctx, req.Email)

	// Generate OTP
	otp := utils.GenerateOTP(5)

	err = h.authUC.SaveOTP(ctx, user.Email, otp, constants.PurposeRegister)

	if err != nil {
		return xerror.BadRequest("Can not save otp")
	}

	// Send OTP via Gmail
	err = h.mailer.SendOTP(req.Email, user.FullName, otp)
	if err != nil {
		return xerror.BadRequest("Can not send otp to email")
	}
	return c.Status(200).JSON("OK")
}
