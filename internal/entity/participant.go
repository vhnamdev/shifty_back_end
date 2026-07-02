package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Participant struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ConversationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_participant_conversation_author" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	AuthorID       uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_participant_conversation_author" json:"author_id"`
	User           User         `gorm:"foreignKey:AuthorID" json:"user,omitempty"`
	IsDeleted      bool         `gorm:"default:false" json:"is_deleted"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at"`
}

func (p *Participant) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
