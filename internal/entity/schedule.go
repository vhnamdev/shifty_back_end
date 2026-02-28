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
	NumberOfMembers int        `gorm:"type:int;default:0;not null" json:"number_of_members"`
	NumberOfShifts  int        `gorm:"type:int;default:0;not null" json:"number_of_shifts"`
	RestaurantID    uuid.UUID  `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant      Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	IsDeleted       bool       `gorm:"default:false" json:"is_deleted"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time  `json:"deleted_at"`
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
