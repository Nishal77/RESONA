package comments

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository { return &Repository{} }

func (r *Repository) Create(ctx context.Context, comment *models.Comment) error {
	return database.DB.WithContext(ctx).Create(comment).Error
}

func (r *Repository) GetByPost(ctx context.Context, postID uuid.UUID, page, limit int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64
	database.DB.WithContext(ctx).Model(&models.Comment{}).
		Where("post_id = ? AND parent_comment_id IS NULL", postID).Count(&total)
	err := database.DB.WithContext(ctx).
		Preload("User").
		Preload("Replies.User").
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Order("like_count DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&comments).Error
	return comments, total, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := database.DB.WithContext(ctx).First(&comment, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (r *Repository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := database.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("comment not found or unauthorized")
	}
	return nil
}

func (r *Repository) IncrementLike(ctx context.Context, id uuid.UUID) error {
	return database.DB.WithContext(ctx).Model(&models.Comment{}).
		Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}
