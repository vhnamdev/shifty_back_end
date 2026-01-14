package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Restaurant struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name  string    `gorm:"type:varchar(100);not null" json:"name"`
	Email string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone string    `gorm:"type:varchar(20)" json:"phone"`
	Address string `gorm:"type:varchar(150)" json:"address"`
	Status  bool   `gorm:"default:false" json:"status"`
	Avatar string `gorm:"type:varchar(255);default:'https://static.vecteezy.com/system/resources/thumbnails/009/292/244/small/default-avatar-icon-of-social-media-user-vector.jpg'" json:"avatar"`
	Users     []User     `gorm:"foreignKey:RestaurantID" json:"users,omitempty"`
	Positions []Position `gorm:"foreignKey:RestaurantID" json:"positions,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
