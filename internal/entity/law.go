package entity

import (
	"shifty-backend/pkg/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Law struct {
	ID            uuid.UUID             `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name          string                `gorm:"type:text;not null" json:"name"`
	Description   string                `gorm:"type:text;not null" json:"description"`
	SeverityLevel constants.SevertyType `gorm:"type:varchar(20);not null" json:"severity_level"`
	RestaurantID  uuid.UUID             `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant    Restaurant            `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	CreatedAt     time.Time             `gorm:"autoCreateTime" json:"created_at"`
}

func (l *Law) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New()
	return
}
