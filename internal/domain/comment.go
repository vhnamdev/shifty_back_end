package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	ImageUrl  *string    `gorm:"type:varchar(255)" json:"image_url"`
	PostID    uuid.UUID  `gorm:"type:uuid;not null" json:"post_id"`
	Post      Post       `gorm:"foreignKey:PostID" json:"post,omitempty"`
	AuthorID  uuid.UUID  `gorm:"type:uuid;not null" json:"author_id"`
	User      User       `gorm:"foreignKey:AuthorID" json:"user,omitempty"`
	ParentID  *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	Replies   []Comment  `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"replies,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
