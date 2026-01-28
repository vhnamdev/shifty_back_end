package entity

import "time"


type UserCache struct {
    UserID       string    `json:"user_id"`
    UserName     string    `json:"user_name"`
    Avatar       string    `json:"avatar"`
    Role         string    `json:"role"`
    Email        string    `json:"email"`
    PhoneNumber  string    `json:"phone_number"` 
    Address      string    `json:"address"`
    PositionID   string    `json:"position_id"`
    RestaurantID string    `json:"restaurant_id"`
    CreatedAt    time.Time `json:"created_at"`
}