package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Engagement struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_post_user_type" json:"post_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_post_user_type" json:"user_id"`
	Type      string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_post_user_type" json:"type"` // like | share | save | view
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (e *Engagement) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
