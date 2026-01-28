package entity

import (
	"time" // Đã thêm import time

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Position struct {
	ID               uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name             string             `gorm:"type:varchar(100);not null" json:"name"`
	Description      string             `gorm:"type:text;not null" json:"description"`
	Salary           int64              `gorm:"type:bigint" json:"salary"`
	RestaurantID     uuid.UUID          `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant       Restaurant         `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	User             []User             `gorm:"foreignKey:PositionID" json:"user_id,omitempty"`
	ShiftRequirement []ShiftRequirement `gorm:"foreignKey:PositionID" json:"shift_requirement,omitempty"`
	ShiftRequest     []ShiftRequest     `gorm:"foreignKey:PositionID" json:"shift_request,omitempty"`
	ShiftAssignment  []ShiftAssignment  `gorm:"foreignKey:PositionID" json:"shift_assignment,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

func (p *Position) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
