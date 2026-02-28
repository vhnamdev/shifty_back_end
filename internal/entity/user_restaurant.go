package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRestaurant struct {
	UserID       uuid.UUID  `gorm:"type:uuid;primaryKey"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;primaryKey"`
	PositionID   *uuid.UUID `gorm:"type:uuid"`
	Position     Position   `gorm:"foreignKey:PositionID"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID"`
	IsBanned     bool       `gorm:"default:false" json:"is_banned"`
	IsDeleted    bool       `gorm:"default:false" json:"is_deleted"`
	DeletedAt    *time.Time `json:"deleted_at"`
	BannedAt     *time.Time `json:"banned_at"`
	JoinedAt     time.Time  `gorm:"autoCreateTime"`
}
