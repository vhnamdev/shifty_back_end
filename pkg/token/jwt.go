package token

import (
	"shifty-backend/pkg/xerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenMaster struct {
	accessSecret    string
	refreshSecret   string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

// Constructor Function
func NewToken(accessSecret, refreshSecret string, accessDuration, refreshDuration time.Duration) *TokenMaster {
	return &TokenMaster{
		accessSecret:    accessSecret,
		refreshSecret:   refreshSecret,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

type UserClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// create new access token function
func (m *TokenMaster) GenerateAccessToken(userId string, role string) (string, error) {
	claims := UserClaims{
		UserID: userId, // Add userid and role to create access token
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.accessSecret))
}

// create new refresh token function
func (m *TokenMaster) GenerateRefreshToken(userID string, role string) (string, error) {
	claims := UserClaims{
		UserID: userID, // Add userid and role to create refresh token
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.refreshSecret))
}

// Verify Access Token
func (m *TokenMaster) VerifyAccessToken(tokenString string) (*UserClaims, error) {
	return m.parseToken(tokenString, m.accessSecret)
}

// Verify Refresh Token
func (m *TokenMaster) VerifyRefreshToken(tokenString string) (*UserClaims, error) {
	return m.parseToken(tokenString, m.refreshSecret)
}

// parseToken is a private helper to validate whether the token is valid using a specific secret key
func (m *TokenMaster) parseToken(tokenString string, secretKey string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, xerror.New(401, "Unexpected signing method")
		}
		return []byte(secretKey), nil

	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, xerror.New(401, "Invalid token")
}
