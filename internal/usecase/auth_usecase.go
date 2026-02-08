package usecase

import (
	"context"
	"shifty-backend/configs"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/mailer"
	"shifty-backend/pkg/token"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"time"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type AuthUseCase interface {
	RegisterLocal(ctx context.Context, user *entity.User) error
	LoginLocal(ctx context.Context, email string, password string, userAgent, clientIP string) (string, string, *entity.User, error)
	LoginGoogle(ctx context.Context, code string, userAgent, clientIP string) (string, string, *entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	SaveOTP(ctx context.Context, email, otp string, purpose constants.OTPPurpose) error
	SendOTP(ctx context.Context, email string, purpose string) error
	ResetPassword(ctx context.Context, email, password string) error
}

type authUseCase struct {
	userRepo       repository.UserRepository
	tokenMaster    *token.TokenMaster
	contextTimeout time.Duration
	redisRepo      repository.RedisRepository
	emailService   *mailer.EmailService
	googleConfig   *configs.GoogleConfig
}

func NewAuthUseCase(repo repository.UserRepository, tokenMaster *token.TokenMaster, timeout time.Duration, redisRepo repository.RedisRepository, emailService *mailer.EmailService, googleConfig *configs.GoogleConfig) AuthUseCase {
	return &authUseCase{
		userRepo:       repo,
		tokenMaster:    tokenMaster,
		contextTimeout: timeout,
		redisRepo:      redisRepo,
		emailService:   emailService,
		googleConfig:   googleConfig,
	}
}

// Register with password and email
func (u *authUseCase) RegisterLocal(ctx context.Context, user *entity.User) error {
	return u.userRepo.Create(ctx, user) // Send ctx and user data to repository
}

// Login with email and password
func (u *authUseCase) LoginLocal(ctx context.Context, email string, password string, userAgent, clientIP string) (string, string, *entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return "", "", nil, xerror.NotFound("User not found")
		}
		return "", "", nil, xerror.Internal("Database error")
	}

	if user.AccountType != "Local" || user.GoogleID == "" {
		return "", "", nil, xerror.BadRequest("Please login via " + user.AccountType)
	}

	// Compare hashpassword in database and password receive from FE
	err = utils.CompareHashAndPassword(password, user.Password)
	if err != nil {
		return "", "", nil, xerror.Unauthorized("Invalid password")
	}

	// Generate access token
	accessToken, err := u.tokenMaster.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return "", "", nil, xerror.Internal("Failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := u.tokenMaster.GenerateRefreshToken(user.ID.String(), user.Role)
	if err != nil {
		return "", "", nil, xerror.Internal("Failed to generate refresh token")
	}

	// Create session
	err = u.redisRepo.CreateSession(ctx, &entity.Session{
		RefreshToken: refreshToken,
		UserID:       user.ID.String(),
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		return "", "", nil, xerror.BadRequest("Can not create session")
	}

	// Save user cache
	err = u.redisRepo.SaveUserCache(ctx, &entity.UserCache{
		UserID:       user.ID.String(),
		UserName:     user.FullName,
		Avatar:       user.Avatar,
		Role:         user.Role,
		Email:        user.Email,
		PhoneNumber:  utils.GetString(user.PhoneNumber),
		Address:      utils.GetString(user.Address),
		PositionID:   user.PositionID.String(),
		RestaurantID: user.RestaurantID.String(),
		CreatedAt:    time.Now(),
	})

	if err != nil {
		return "", "", nil, xerror.BadRequest("Can not cache user data")
	}

	return accessToken, refreshToken, user, nil
}

// Login By Google
func (u *authUseCase) LoginGoogle(ctx context.Context, code string, userAgent, clientIP string) (string, string, *entity.User, error) {

	// 
	token, err := u.googleConfig.GoogleLoginConfig.Exchange(ctx, code)
	if err != nil {
		return "", "", nil, xerror.Unauthorized("Failed to exchange token with Google")
	}

	oauth2Service, err := oauth2.NewService(ctx, option.WithTokenSource(u.googleConfig.GoogleLoginConfig.TokenSource(ctx, token)))

	if err != nil {
		return "", "", nil, xerror.Internal("Failed to create Google service")
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return "", "", nil, xerror.Internal("Failed to get user into from Google")
	}

	user, err := u.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			fullName := userInfo.FamilyName + userInfo.GivenName
			hashPassword, _ := utils.HashPassword("Google password")
			newUser := &entity.User{
				FullName:    fullName,
				Email:       userInfo.Email,
				Avatar:      userInfo.Picture,
				AccountType: "Google",
				Role:        "User",
				GoogleID:    userInfo.Id,
				Status:      true,
				Password:    hashPassword,
			}
			if err := u.userRepo.Create(ctx, newUser); err != nil {
				return "", "", nil, xerror.Internal("Failed to create User")
			}
			user = newUser
		} else {
			return "", "", nil, xerror.Internal("Database error")
		}
	}
	accessToken, err := u.tokenMaster.GenerateAccessToken(user.ID.String(), user.Role)

	if err != nil {
		return "", "", nil, xerror.Internal("Failed to generate access token")
	}

	refreshToken, err := u.tokenMaster.GenerateRefreshToken(user.ID.String(), user.Role)

	if err != nil {
		return "", "", nil, xerror.Internal("Failed to generate refresh token")
	}

	// Create session
	err = u.redisRepo.CreateSession(ctx, &entity.Session{
		RefreshToken: refreshToken,
		UserID:       user.ID.String(),
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		return "", "", nil, xerror.BadRequest("Can not create session")
	}

	return accessToken, refreshToken, user, nil

}

// Find user by email
func (u *authUseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	// Get user
	user, err := u.userRepo.GetByEmail(ctx, email)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return nil, xerror.NotFound("User not found")
		}
		return nil, xerror.Internal("Database error")
	}

	return user, nil
}

// Save otp into redis
func (u *authUseCase) SaveOTP(ctx context.Context, email, otp string, purpose constants.OTPPurpose) error {

	// Save OTP
	err := u.redisRepo.SaveOTP(ctx, email, otp, purpose)
	if err != nil {
		return err
	}
	return nil

}

// Func Send OTP to register
func (u *authUseCase) SendOTP(ctx context.Context, email string, purpose string) error {

	// Get otp purpose
	otpPurpose := constants.OTPPurpose(purpose)

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, email)

	// Create user's name variable
	var userName string

	// Use switch case to check the purpose of OTP
	switch otpPurpose {

	// If user want to regist account
	case constants.PurposeRegister:

		// Check if user is exist
		if user != nil {
			return xerror.BadRequest("Email is already exist")
		}
		// Return error if it's a system failure
		if !utils.IsRecordNotFoundError(err) {
			return xerror.Internal("Database error")
		}
		// Set user name
		userName = "New User"
	case constants.PurposeResetPassword:

		// If err or user not found
		if err != nil {
			if utils.IsRecordNotFoundError(err) {
				return xerror.Internal("User not found")
			}
			return xerror.Internal("Database error")
		}
		// Set username == user's name
		userName = user.FullName
	default:
		return xerror.BadRequest("Invalid OTP purpose")
	}

	// Generate OTP with 5 digits
	otp := utils.GenerateOTP(5)

	// Save OTP into Redis database
	err = u.redisRepo.SaveOTP(ctx, email, otp, otpPurpose)

	if err != nil {
		return xerror.Internal("Failed to save OTP")
	}

	go func() {
		_ = u.emailService.SendOTP(email, userName, otp)
	}()
	return nil
}

// Reset Password
func (u *authUseCase) ResetPassword(ctx context.Context, email, password string) error {

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, email)

	if err != nil {
		if utils.IsRecordNotFoundError(err) {
			return xerror.NotFound("User not found")
		}
		return xerror.Internal("Database error")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return xerror.BadRequest("Failed to hash password")
	}

	// Replace old password by new hashed passowrd
	user.Password = hashedPassword

	// Update password
	if err := u.userRepo.Update(ctx, user); err != nil {
		return xerror.Internal("Failed to update password")
	}
	return nil
}
