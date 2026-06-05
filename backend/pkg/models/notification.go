package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_notif_user_read;index:idx_notif_user_created" json:"user_id"`
	Type      string     `gorm:"type:varchar(20);not null" json:"type"` // like | comment | follow | trending
	ActorID   *uuid.UUID `gorm:"type:uuid" json:"actor_id"`
	Actor     *User      `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
	PostID    *uuid.UUID `gorm:"type:uuid" json:"post_id"`
	Message   *string    `gorm:"type:text" json:"message"`
	Read      bool       `gorm:"not null;default:false;index:idx_notif_user_read" json:"read"`
	CreatedAt time.Time  `gorm:"index:idx_notif_user_created,sort:desc" json:"created_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
