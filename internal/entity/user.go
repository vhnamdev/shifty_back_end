package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FullName          string      `gorm:"type:varchar(100);not null" json:"full_name"`
	Email             string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password          string      `gorm:"not null" json:"-"`
	PhoneNumber       *string     `gorm:"type:varchar(20)" json:"phone_number"`
	Role              string      `gorm:"type:varchar(20);default:'user'" json:"role"`
	Address           *string     `gorm:"type:varchar(150)" json:"address"`
	Status            bool        `gorm:"default:false" json:"status"`
	GoogleID          string      `gorm:"index"`
	AccountType       string      `gorm:"default:'Local'" json:"account_type"`
	Avatar            string      `gorm:"type:varchar(255);default:'https://static.vecteezy.com/system/resources/thumbnails/009/292/244/small/default-avatar-icon-of-social-media-user-vector.jpg'" json:"avatar"`
	PositionID        *uuid.UUID  `gorm:"type:uuid" json:"position_id,omitempty"`
	Position          *Position   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	RestaurantID      *uuid.UUID  `gorm:"type:uuid" json:"restaurant_id,omitempty"`
	Restaurant        *Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	ReceivedFeedbacks []Feedback  `gorm:"foreignKey:MemberID" json:"received_feedbacks,omitempty"`
	GivenFeedbacks    []Feedback  `gorm:"foreignKey:ReviewerID" json:"given_feedbacks,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
