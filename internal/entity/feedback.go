package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	MemberID   uuid.UUID `gorm:"type:uuid;not null" json:"member_id,omitempty"`
	Member     User      `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	ReviewerID uuid.UUID `gorm:"type:uuid;not null" json:"reviewer_id,omitempty"`
	Reviewer   User      `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
