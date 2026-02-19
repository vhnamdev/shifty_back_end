package entity

import (
	"shifty-backend/pkg/constants"
	"time"
)

type UserCache struct {
	UserID       string             `json:"user_id"`
	UserName     string             `json:"user_name"`
	Avatar       string             `json:"avatar"`
	Role         constants.UserRole `json:"role"`
	Email        string             `json:"email"`
	PhoneNumber  string             `json:"phone_number"`
	Address      string             `json:"address"`
	CreatedAt    time.Time          `json:"created_at"`
}
