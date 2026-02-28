package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Restaurant struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name            string           `gorm:"type:varchar(100);not null" json:"name"`
	Email           string           `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone           string           `gorm:"type:varchar(20)" json:"phone"`
	Address         string           `gorm:"type:varchar(150)" json:"address"`
	Status          bool             `gorm:"default:true" json:"status"`
	IsDeleted       bool             `gorm:"default:false" json:"is_deleted"`
	Avatar          string           `gorm:"type:varchar(255)" json:"avatar"`
	UserRestaurants []UserRestaurant `gorm:"foreignKey:RestaurantID" json:"user_restaurants,omitempty"`
	Positions       []Position       `gorm:"foreignKey:RestaurantID" json:"positions,omitempty"`
	Laws            []Law            `gorm:"foreignKey:RestaurantID" json:"laws,omitempty"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       time.Time        `json:"deleted_at"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
