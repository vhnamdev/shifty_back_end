package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reaction struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_reaction_post_author" json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	AuthorID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_reaction_post_author" json:"author_id"`
	User      User      `gorm:"foreignKey:AuthorID" json:"user,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Reaction) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
