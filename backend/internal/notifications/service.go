package notifications

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Create(ctx context.Context, userID, actorID uuid.UUID, notifType string, postID *uuid.UUID, msg *string) {
	n := &models.Notification{
		UserID:  userID,
		ActorID: &actorID,
		Type:    notifType,
		PostID:  postID,
		Message: msg,
	}
	database.DB.WithContext(ctx).Create(n)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Notification, int64, error) {
	var notifs []models.Notification
	var total int64
	database.DB.WithContext(ctx).Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total)
	err := database.DB.WithContext(ctx).
		Preload("Actor").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&notifs).Error
	return notifs, total, err
}

func (s *Service) UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND read = false", userID).Count(&count).Error
	return count, err
}

func (s *Service) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	result := database.DB.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).Update("read", true)
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}
	return result.Error
}

func (s *Service) MarkAllRead(ctx context.Context, userID uuid.UUID) (int64, error) {
	result := database.DB.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND read = false", userID).Update("read", true)
	return result.RowsAffected, result.Error
}

// unused
var _ = gorm.ErrRecordNotFound
