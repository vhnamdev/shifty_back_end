package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ShiftRule struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type         string         `gorm:"type:varchar(50); not null;index" json:"type"`
	Name         string         `gorm:"type:varchar(100); not null" json:"name"`
	Config       datatypes.JSON `gorm:"type:jsonb" json:"config"`
	IsDeleted    bool           `gorm:"default:false" json:"is_deleted"`
	RestaurantID uuid.UUID      `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant   Restaurant     `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at"`
}

func (s *ShiftRule) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
