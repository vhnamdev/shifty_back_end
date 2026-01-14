package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Content string    `gorm:"type:text;not null" json:"content"`
	ImageUrl string `gorm:"type:varchar(255)" json:"image_url"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	AuthorID     uuid.UUID  `gorm:"type:uuid;not null" json:"author_id"`
	User         User       `gorm:"foreignKey:AuthorID" json:"user,omitempty"`
	Reactions []Reaction `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"reactions,omitempty"`
	Comments  []Comment  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
