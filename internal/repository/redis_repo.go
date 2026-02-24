package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"shifty-backend/internal/dto"
	"shifty-backend/internal/entity"
	"shifty-backend/pkg/constants"
	"shifty-backend/pkg/xerror"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	// Session
	CreateSession(ctx context.Context, session *entity.Session) error
	GetSession(ctx context.Context, refreshToken string) (*entity.Session, error)
	DeleteSession(ctx context.Context, userID string, refreshToken string) error
	DeleteAllSessions(ctx context.Context, userID string) error

	// User Cache
	SaveUserCache(ctx context.Context, user *entity.UserCache) error
	GetUserCache(ctx context.Context, userID string) (*entity.UserCache, error)
	DeleteUserCache(ctx context.Context, userID string) error

	// OTP
	SaveOTP(ctx context.Context, email string, otp string, purpose constants.OTPPurpose) error
	VerifyOTP(ctx context.Context, email string, inputOTP string, purpose constants.OTPPurpose) error

	// Online
	SetUserStatus(ctx context.Context, userID string, isOnline bool) error
	GetUserStatus(ctx context.Context, userID string) (bool, error)

	// Invite Code
	SaveInviteCode(ctx context.Context, inviteCode, email, positonID, resID string) error
	VerifyInviteCode(ctx context.Context, email, inviteCode string) (*dto.InviteData, error)
}

type redisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) RedisRepository {
	return &redisRepo{client: client}
}

// Create session when users login
func (r *redisRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	dataBytes, err := json.Marshal(session) // Serialize Entity -> JSON bytes
	if err != nil {
		return xerror.Internal("Failed to encode session data")
	}
	pipe := r.client.Pipeline() // Use Pipeline to batch commands (reduce network round-trips)

	sessionKey := fmt.Sprintf("session:%s", session.RefreshToken) // Set main session key with TTL
	pipe.Set(ctx, sessionKey, dataBytes, 30*24*time.Hour)         //  TTL is 30 days

	userSessionKey := fmt.Sprintf("user:session:%s", session.UserID) // Add to user's device list
	pipe.SAdd(ctx, userSessionKey, session.RefreshToken)

	_, err = pipe.Exec(ctx) //Execute all commands
	return err
}

// GetSession retrieves session data by refresh token
func (r *redisRepo) GetSession(ctx context.Context, refreshToken string) (*entity.Session, error) {
	// make a key with refreshtoken
	key := fmt.Sprintf("session:%s", refreshToken)

	// Check if key exists and get value
	val, err := r.client.Get(ctx, key).Result()

	// Key not found (Expired or Invalid token)
	if err == redis.Nil {
		return nil, nil
	}

	// if redis crack or disconnect network
	if err != nil {
		return nil, err
	}

	var session entity.Session // create session variable use to convert data json in redis to struct session in go struct
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, xerror.Internal("Fail to parse session data")
	}
	return &session, nil // return session with error is nil
}

// Delete one session with userID and refreshToken
func (r *redisRepo) DeleteSession(ctx context.Context, userID string, refreshToken string) error {
	// create pipe
	pipe := r.client.Pipeline()

	// create key session
	key := fmt.Sprintf("session:%s", refreshToken)

	// add delete key into pipe
	pipe.Del(ctx, key)

	// create user session key
	userSessionKey := fmt.Sprintf("user:session:%s", userID)

	// delete many token in user session
	pipe.SRem(ctx, userSessionKey, refreshToken)

	// excute commands
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteAllSessions invalidates ALL active sessions for a user (Logout All)
func (r *redisRepo) DeleteAllSessions(ctx context.Context, userID string) error {

	// Construct the key for user's session list
	userSessionKey := fmt.Sprintf("user:session:%s", userID)

	// Retrieve all refresh tokens belonging to this user
	tokens, err := r.client.SMembers(ctx, userSessionKey).Result()

	if err != nil {
		return err
	}

	//No sessions found -> nothing to do
	if len(tokens) == 0 {
		return nil
	}

	// Init Pipeline for atomic/batch execution
	pipe := r.client.Pipeline()

	// Queue commands: Delete each individual session data
	for _, token := range tokens {
		sessionKey := fmt.Sprintf("session:%s", token)
		pipe.Del(ctx, sessionKey)
	}

	// Queue command: Delete the user's session list (the Set)
	pipe.Del(ctx, userSessionKey)

	// Execute the transaction
	_, err = pipe.Exec(ctx)
	return err
}

// Save User Cache when login or update user's information
func (r *redisRepo) SaveUserCache(ctx context.Context, user *entity.UserCache) error {

	// Serialize Entity -> JSON bytes
	dataBytes, err := json.Marshal(user)

	if err != nil {
		return xerror.Internal("Failed to encode user data")
	}

	// Create user cache key
	userCacheKey := fmt.Sprintf("cache:user:%s", user.UserID)

	// excecute command save user data into redis
	_, err = r.client.Set(ctx, userCacheKey, dataBytes, 30*24*time.Hour).Result()

	return err
}

// Get user's cache by userID
func (r *redisRepo) GetUserCache(ctx context.Context, userID string) (*entity.UserCache, error) {

	// Create user's cache key
	userCacheKey := fmt.Sprintf("cache:user:%s", userID)

	// Get value with key
	val, err := r.client.Get(ctx, userCacheKey).Result()

	// If data not found return nil and nil
	if err == redis.Nil {
		return nil, nil
	}

	// if err return err
	if err != nil {
		return nil, err
	}

	// create user with type is entity User
	var user entity.UserCache

	// Convert data from json to go struct and return
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, xerror.Internal("Fail to parse user cache")
	}
	return &user, nil
}

func (r *redisRepo) DeleteUserCache(ctx context.Context, userID string) error {
	userCacheKey := fmt.Sprintf("cache:user:%s", userID)
	return r.client.Del(ctx, userCacheKey).Err()
}

// Save OTP into redis by email and purpose
func (r *redisRepo) SaveOTP(ctx context.Context, email string, otp string, purpose constants.OTPPurpose) error {
	// Create otp key
	otpKey := fmt.Sprintf("otp:%s:%s", purpose, email)

	// Set otp with key into redis
	return r.client.Set(ctx, otpKey, otp, 5*time.Minute).Err()
}

// Verify OTP
func (r *redisRepo) VerifyOTP(ctx context.Context, email string, inputOTP string, purpose constants.OTPPurpose) error {

	// Create otp key
	otpKey := fmt.Sprintf("otp:%s:%s", purpose, email)

	// Get otp from key
	val, err := r.client.Get(ctx, otpKey).Result()

	// If value is nil return error
	if err == redis.Nil {
		return xerror.BadRequest("OTP expired or invalid")
	}

	// If err return
	if err != nil {
		return err
	}

	// If otp saved in db is not match with input OTP return error
	if val != inputOTP {
		return xerror.BadRequest("Wrong OTP code")
	}

	// Delete otp key
	// Ignore delete error. Even if cleanup fails, the user has successfully verified the OTP.
	// Returning an error here would confuse the user.
	r.client.Del(ctx, otpKey)
	return nil
}

// SetUserStatus updates user's online status.
// isOnline=true for login/active, isOnline=false for logout.
func (r *redisRepo) SetUserStatus(ctx context.Context, userID string, isOnline bool) error {

	// create user status key
	userStatusKey := fmt.Sprintf("user:online:%s", userID)

	// Prepare data map
	data := map[string]interface{}{
		"status":     isOnline,
		"time_stamp": time.Now().Unix(),
	}

	// Save to Redis Hash
	if err := r.client.HSet(ctx, userStatusKey, data).Err(); err != nil {
		return err
	}

	// If Online: Set TTL (Heartbeat).
	// If user is inactive for 5 minutes, the key expires (Status -> Offline)
	if isOnline {
		r.client.Expire(ctx, userStatusKey, 5*time.Minute)
	}
	return nil
}

// GetUserStatus retrieves the current online status of a user.
func (r *redisRepo) GetUserStatus(ctx context.Context, userID string) (bool, error) {

	// Create usr status key
	userStatusKey := fmt.Sprintf("user:online:%s", userID)

	// Fetch only the 'status' field to save bandwidth
	val, err := r.client.HGet(ctx, userStatusKey, "status").Result()

	// Key doesn't exist => User is Offline
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	// Parse Redis string result to bool
	isOnline, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return isOnline, nil
}

// Save Invite Code
func (r *redisRepo) SaveInviteCode(ctx context.Context, inviteCode, email, positonID, resID string) error {

	// Create key
	inviteCodeKey := fmt.Sprintf("invite:%s", email)

	// Put invite code and position id into struct
	inviteData := dto.InviteData{
		InviteCode:   inviteCode,
		PositionID:   positonID,
		RestaurantID: resID,
	}

	// Convert invited data to json to save with key
	dataBytes, err := json.Marshal(inviteData)

	if err != nil {
		return xerror.Internal("Failed to encode invite data")
	}

	return r.client.Set(ctx, inviteCodeKey, dataBytes, 24*time.Hour).Err()
}

func (r *redisRepo) VerifyInviteCode(ctx context.Context, email, inviteCode string) (*dto.InviteData, error) {
	inviteCodeKey := fmt.Sprintf("invite:%s", email)

	val, err := r.client.Get(ctx, inviteCodeKey).Result()

	if err == redis.Nil {
		return nil, xerror.BadRequest("Invite code expired or invalid")
	}

	if err != nil {
		return nil, err
	}

	var inviteData dto.InviteData
	if err := json.Unmarshal([]byte(val), &inviteData); err != nil {
		return nil, xerror.Internal("Can not read invite data")
	}

	if inviteData.InviteCode != inviteCode {
		return nil, xerror.BadRequest("Incorrect invitation code")
	}

	_ = r.client.Del(ctx, inviteCodeKey).Err()

	return &inviteData, nil

}
