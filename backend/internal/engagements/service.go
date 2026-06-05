package engagements

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Record(ctx context.Context, postID, userID uuid.UUID, engType string) error {
	eng := &models.Engagement{PostID: postID, UserID: userID, Type: engType}
	if err := database.DB.WithContext(ctx).Create(eng).Error; err != nil {
		return fmt.Errorf("duplicate engagement")
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, postID, userID uuid.UUID, engType string) error {
	return database.DB.WithContext(ctx).
		Where("post_id = ? AND user_id = ? AND type = ?", postID, userID, engType).
		Delete(&models.Engagement{}).Error
}

func (s *Service) HasEngaged(ctx context.Context, postID, userID uuid.UUID, engType string) (bool, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.Engagement{}).
		Where("post_id = ? AND user_id = ? AND type = ?", postID, userID, engType).
		Count(&count).Error
	return count > 0, err
}
