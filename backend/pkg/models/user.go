package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username            string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email               string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash        *string    `gorm:"type:varchar" json:"-"`
	GoogleID            *string    `gorm:"type:varchar;uniqueIndex" json:"-"`
	FullName            *string    `gorm:"type:varchar(100)" json:"full_name"`
	AvatarURL           *string    `gorm:"type:varchar" json:"avatar_url"`
	Bio                 *string    `gorm:"type:text" json:"bio"`
	PrimaryLanguage     string     `gorm:"type:varchar(20);not null;default:'kannada'" json:"primary_language"`
	State               *string    `gorm:"type:varchar(100)" json:"state"`
	City                *string    `gorm:"type:varchar(100)" json:"city"`
	VRSTotal            float64    `gorm:"default:0" json:"vrs_total"`
	FollowerCount       int        `gorm:"default:0" json:"follower_count"`
	FollowingCount      int        `gorm:"default:0" json:"following_count"`
	OnboardingCompleted bool       `gorm:"default:false" json:"onboarding_completed"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
