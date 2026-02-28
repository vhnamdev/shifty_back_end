package entity

import (
	"shifty-backend/pkg/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FullName          string             `gorm:"type:varchar(100);not null" json:"full_name"`
	Email             string             `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password          string             `gorm:"not null" json:"-"`
	PhoneNumber       *string            `gorm:"type:varchar(20)" json:"phone_number"`
	Role              constants.UserRole `gorm:"type:varchar(20);default:'USER'" json:"role"`
	Address           *string            `gorm:"type:varchar(150)" json:"address"`
	Status            bool               `gorm:"default:true" json:"status"`
	IsDeleted         bool               `gorm:"default:false" json:"is_deleted"`
	GoogleID          string             `gorm:"index"`
	AccountType       string             `gorm:"default:'Local'" json:"account_type"`
	Avatar            string             `gorm:"type:varchar(255)" json:"avatar"`
	UserRestaurants   []UserRestaurant   `gorm:"foreignKey:UserID" json:"user_restaurants,omitempty"`
	ReceivedFeedbacks []Feedback         `gorm:"foreignKey:MemberID" json:"received_feedbacks,omitempty"`
	GivenFeedbacks    []Feedback         `gorm:"foreignKey:ReviewerID" json:"given_feedbacks,omitempty"`
	CreatedAt         time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	DeletedAt         *time.Time         `json:"deletedAt" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
