package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Content      string     `gorm:"type:text;not null" json:"content"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	MemberID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"member_id,omitempty"`
	Member       User       `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	ReviewerID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"reviewer_id,omitempty"`
	Reviewer     User       `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	IsDeleted    bool       `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

func (f *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	return
}
