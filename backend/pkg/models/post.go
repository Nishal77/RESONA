package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID                uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	User                  User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ContentText           *string    `gorm:"type:text" json:"content_text"`
	MediaURL              *string    `gorm:"type:varchar" json:"media_url"`
	MediaType             *string    `gorm:"type:varchar(10)" json:"media_type"` // image | video | null
	CommunityID           *uuid.UUID `gorm:"type:uuid;index" json:"community_id"`
	Community             *Community `gorm:"foreignKey:CommunityID" json:"community,omitempty"`
	DetectedLanguage      *string    `gorm:"type:varchar(20)" json:"detected_language"`
	LanguageConfidence    *float64   `json:"language_confidence"`
	LanguageLocalityScore float64    `gorm:"not null;default:0.3" json:"language_locality_score"`
	VRSScore              float64    `gorm:"not null;default:0;index:idx_posts_vrs_score,sort:desc" json:"vrs_score"`
	LikeCount             int        `gorm:"not null;default:0" json:"like_count"`
	CommentCount          int        `gorm:"not null;default:0" json:"comment_count"`
	ShareCount            int        `gorm:"not null;default:0" json:"share_count"`
	ViewCount             int        `gorm:"not null;default:0" json:"view_count"`
	SaveCount             int        `gorm:"not null;default:0" json:"save_count"`
	Tags                  []Tag      `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
