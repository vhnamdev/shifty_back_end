package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reaction struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	AuthorID  uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	User      User      `gorm:"foreignKey:AuthorID" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Reaction) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
