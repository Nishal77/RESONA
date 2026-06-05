package posts

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

func (r *Repository) Create(ctx context.Context, post *models.Post) error {
	return database.DB.WithContext(ctx).Create(post).Error
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := database.DB.WithContext(ctx).
		Preload("User").Preload("Community").Preload("Tags").
		First(&post, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &post, err
}

func (r *Repository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := database.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("post not found or unauthorized")
	}
	return nil
}

func (r *Repository) GetFeed(ctx context.Context, language string, page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	q := database.DB.WithContext(ctx).Model(&models.Post{})
	if language != "" {
		q = q.Where("detected_language = ?", language)
	}
	q.Count(&total)
	err := q.Preload("User").Preload("Community").Preload("Tags").
		Order("vrs_score DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&posts).Error
	return posts, total, err
}

func (r *Repository) GetPersonalizedFeed(ctx context.Context, userID uuid.UUID, followedIDs, communityIDs []uuid.UUID, language string, page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	q := database.DB.WithContext(ctx).Model(&models.Post{}).
		Where("user_id IN (?) OR community_id IN (?)", followedIDs, communityIDs)
	if language != "" {
		q = q.Where("detected_language = ?", language)
	}
	q.Count(&total)
	err := q.Preload("User").Preload("Community").Preload("Tags").
		Order("vrs_score DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&posts).Error
	return posts, total, err
}

func (r *Repository) Update(ctx context.Context, post *models.Post) error {
	return database.DB.WithContext(ctx).Save(post).Error
}

func (r *Repository) IncrementCount(ctx context.Context, postID uuid.UUID, column string, delta int) error {
	return database.DB.WithContext(ctx).Model(&models.Post{}).
		Where("id = ?", postID).
		UpdateColumn(column, gorm.Expr(column+" + ?", delta)).Error
}

func (r *Repository) DecrementCount(ctx context.Context, postID uuid.UUID, column string) error {
	return database.DB.WithContext(ctx).Model(&models.Post{}).
		Where("id = ?", postID).
		UpdateColumn(column, gorm.Expr("GREATEST("+column+" - 1, 0)")).Error
}

func (r *Repository) AssociateTags(ctx context.Context, post *models.Post, tags []models.Tag) error {
	return database.DB.WithContext(ctx).Model(post).Association("Tags").Append(tags)
}
