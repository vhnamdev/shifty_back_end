package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shift struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	NumberOfMembers *int      `gorm:"type:int" json:"number_of_members"`
	Type            string    `gorm:"type:varchar(50);index" json:"type"`
	ScheduleID      uuid.UUID `gorm:"type:uuid;not null" json:"schedule_id"`
	Schedule        Schedule  `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	IsHoliday       bool      `gorm:"default:false" json:"is_holiday"`
	WageMultiplier  float64   `gorm:"type:decimal(3,2);default:1.0" json:"wage_multiplier"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (r *Shift) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
