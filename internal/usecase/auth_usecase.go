package usecase

import (
	"context"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/token"
	"shifty-backend/pkg/utils"
	"shifty-backend/pkg/xerror"
	"time"
)

type AuthUseCase interface {
	RegisterLocal(ctx context.Context, user *entity.User) error
	LoginLocal(ctx context.Context, email string, password string, userAgent, clientIP string) (string, string, *entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	SaveOTP(ctx context.Context, email, otp string, purpose constants.OTPPurpose) error
	// RegisterByGoogle(ctx context.Context) error
}

type authUseCase struct {
	userRepo       repository.UserRepository
	tokenMaster    *token.TokenMaster
	contextTimeout time.Duration
	redisRepo      repository.RedisRepository
}

func NewAuthUseCase(repo repository.UserRepository, tokenMaster *token.TokenMaster, timeout time.Duration, redisRepo repository.RedisRepository) AuthUseCase {
	return &authUseCase{
		userRepo:       repo,
		tokenMaster:    tokenMaster,
		contextTimeout: timeout,
		redisRepo:      redisRepo,
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
		return "", "", nil, xerror.Internal("Database error")
	}
	if user == nil {
		return "", "", nil, xerror.NotFound("User not found")
	}
	if user.AccountType != "Local" {
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

// Find user by email
func (u *authUseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	// Get user
	user, err := u.userRepo.GetByEmail(ctx, email)

	if user == nil {
		return nil, xerror.NotFound("Can not find user")
	}

	if err != nil {
		return nil, err
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

// func (u *authUseCase )RegisterByGoogle(ctx context.Context) error{

// }
