package entity

import (
	"time" // Đã thêm import time

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Position struct {
	ID                  uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name                string             `gorm:"type:varchar(100);not null" json:"name"`
	Description         string             `gorm:"type:text;not null" json:"description"`
	Rank                int                `gorm:"default:5;not null" json:"rank"`
	CanUpdateRestaurant bool               `gorm:"default:false" json:"can_update_restaurant"`
	CanDeleteRestaurant bool               `gorm:"default:false" json:"can_delete_restaurant"`
	Salary              *int64             `gorm:"type:bigint" json:"salary"`
	IsDeleted           bool               `gorm:"default:false" json:"is_deleted"`
	RestaurantID        uuid.UUID          `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant          Restaurant         `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	UserRestaurants     []UserRestaurant   `gorm:"foreignKey:PositionID" json:"users,omitempty"`
	ShiftRequirements   []ShiftRequirement `gorm:"foreignKey:PositionID" json:"shift_requirement,omitempty"`
	ShiftRequests       []ShiftRequest     `gorm:"foreignKey:PositionID" json:"shift_request,omitempty"`
	ShiftAssignments    []ShiftAssignment  `gorm:"foreignKey:PositionID" json:"shift_assignment,omitempty"`
	CreatedAt           time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	DeletedAt           time.Time          `json:"deleted_at"`
}

func (p *Position) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
