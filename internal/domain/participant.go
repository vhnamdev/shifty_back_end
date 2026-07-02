package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Participant struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ConversationID uuid.UUID    `gorm:"type:uuid;not null" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	AuthorID       uuid.UUID    `gorm:"type:uuid;not null" json:"author_id"`
	User           User         `gorm:"foreignKey:AuthorID" json:"user,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func (p *Participant) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
