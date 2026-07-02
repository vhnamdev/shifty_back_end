package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	ConversationID uuid.UUID    `gorm:"type:uuid;not null;index" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	SenderID       uuid.UUID    `gorm:"type:uuid;not null;index" json:"sender_id"`
	Sender         User         `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Content        string       `gorm:"type:text" json:"content"`
	ImageUrl       *string      `gorm:"type:varchar(255)" json:"image_url"`
	IsDeleted      bool         `gorm:"default:false" json:"is_deleted"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
