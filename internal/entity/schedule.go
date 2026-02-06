package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Schedule struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         time.Time  `json:"end_time"`
	NumberOfMembers *int       `gorm:"type:int;not null" json:"number_of_members"`
	NumberOfShifts  *int       `gorm:"type:int;not null" json:"number_of_shifts"`
	RestaurantID    uuid.UUID  `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant      Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (r *Schedule) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
