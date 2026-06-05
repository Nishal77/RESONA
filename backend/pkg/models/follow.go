package models

import "time"

import "github.com/google/uuid"

type Follow struct {
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;primaryKey;index" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;primaryKey;index" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}
