package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_tag_name_lang" json:"name"`
	Language      *string   `gorm:"type:varchar(20);uniqueIndex:idx_tag_name_lang" json:"language"`
	UsageCount    int       `gorm:"not null;default:0" json:"usage_count"`
	TrendingScore float64   `gorm:"not null;default:0" json:"trending_score"`
	CreatedAt     time.Time `json:"created_at"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
