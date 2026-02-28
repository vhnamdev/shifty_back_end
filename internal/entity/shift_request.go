package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftRequest struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShiftID    uuid.UUID  `gorm:"type:uuid;not null" json:"shift_id"`
	Shift      Shift      `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
	PositionID *uuid.UUID `gorm:"type:uuid" json:"position_id"`
	Position   Position   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Note       *string    `gorm:"type:text" json:"note"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (s *ShiftRequest) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
