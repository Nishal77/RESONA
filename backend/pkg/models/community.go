package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	ID                 uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name               string       `gorm:"type:varchar(100);not null" json:"name"`
	Slug               string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description        *string      `gorm:"type:text" json:"description"`
	PrimaryLanguage    string       `gorm:"type:varchar(20);not null;index" json:"primary_language"`
	SecondaryLanguages StringArray  `gorm:"type:varchar(20)[];default:'{}'" json:"secondary_languages"`
	AvatarURL          *string      `gorm:"type:varchar" json:"avatar_url"`
	BannerURL          *string      `gorm:"type:varchar" json:"banner_url"`
	MemberCount        int          `gorm:"not null;default:0" json:"member_count"`
	PostCount          int          `gorm:"not null;default:0" json:"post_count"`
	SnapOfWeekPostID   *uuid.UUID   `gorm:"type:uuid" json:"snap_of_week_post_id"`
	CreatedBy          *uuid.UUID   `gorm:"type:uuid" json:"created_by"`
	Creator            *User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt          time.Time    `json:"created_at"`
}

func (c *Community) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CommunityMember struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CommunityID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_community_user" json:"community_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_community_user" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role        string    `gorm:"type:varchar(20);not null;default:'member'" json:"role"` // member | moderator | admin
	JoinedAt    time.Time `json:"joined_at"`
}

func (cm *CommunityMember) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == uuid.Nil {
		cm.ID = uuid.New()
	}
	return nil
}
