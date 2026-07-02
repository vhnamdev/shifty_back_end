package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftAssignment struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime time.Time  `json:"check_out_time"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShiftID      uuid.UUID  `gorm:"type:uuid;not null" json:"shift_id"`
	Shift        Shift      `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
	PositionID   *uuid.UUID `gorm:"type:uuid" json:"position_id"`
	Position     Position   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Note         *string    `gorm:"type:text" json:"note"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (r *ShiftAssignment) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
