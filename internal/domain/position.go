package domain

import (
	"time" // Đã thêm import time

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Position struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Salary int64 `gorm:"type:bigint" json:"salary"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Position) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
