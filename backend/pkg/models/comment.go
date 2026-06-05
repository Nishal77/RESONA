package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PostID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"post_id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User            User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content         string     `gorm:"type:text;not null" json:"content"`
	ParentCommentID *uuid.UUID `gorm:"type:uuid;index" json:"parent_comment_id"`
	Replies         []Comment  `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`
	LikeCount       int        `gorm:"not null;default:0" json:"like_count"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
