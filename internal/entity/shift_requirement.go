package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftRequirement struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
	Quantity   *int       `gorm:"type:int;not null" json:"number_of_members"`
	ShiftID    uuid.UUID  `gorm:"type:uuid;not null" json:"shift_id"`
	Shift      Shift      `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
	PositionID *uuid.UUID `gorm:"type:uuid" json:"position_id"`
	Position   Position   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Note       *string    `gorm:"type:text" json:"note"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (r *ShiftRequirement) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
