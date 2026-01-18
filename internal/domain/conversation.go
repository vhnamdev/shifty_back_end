package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ConvTypeDirect = "DIRECT"
	ConvTypeGroup  = "GROUP"
)

type Conversation struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type          string        `gorm:"type:varchar(20);default:'DIRECT';index" json:"type"`
	Name          *string       `gorm:"type:varchar(100)" json:"name,omitempty"`
	Avatar        *string       `gorm:"type:text" json:"image_url,omitempty"`
	LastMessageAt *time.Time    `gorm:"index" json:"last_message_at"`
	Participants  []Participant `gorm:"foreignKey:ConversationID" json:"participants,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (p *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
