package communities

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository { return &Repository{} }

func (r *Repository) Create(ctx context.Context, c *models.Community) error {
	return database.DB.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindBySlug(ctx context.Context, slug string) (*models.Community, error) {
	var c models.Community
	err := database.DB.WithContext(ctx).Preload("Creator").First(&c, "slug = ?", slug).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &c, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Community, error) {
	var c models.Community
	err := database.DB.WithContext(ctx).First(&c, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &c, err
}

func (r *Repository) Update(ctx context.Context, c *models.Community) error {
	return database.DB.WithContext(ctx).Save(c).Error
}

func (r *Repository) List(ctx context.Context, language string, page, limit int) ([]models.Community, int64, error) {
	var communities []models.Community
	var total int64
	q := database.DB.WithContext(ctx).Model(&models.Community{})
	if language != "" {
		q = q.Where("primary_language = ?", language)
	}
	q.Count(&total)
	err := q.Order("member_count DESC").Offset((page - 1) * limit).Limit(limit).Find(&communities).Error
	return communities, total, err
}

func (r *Repository) GetMembers(ctx context.Context, communityID uuid.UUID, page, limit int) ([]models.CommunityMember, int64, error) {
	var members []models.CommunityMember
	var total int64
	database.DB.WithContext(ctx).Model(&models.CommunityMember{}).Where("community_id = ?", communityID).Count(&total)
	err := database.DB.WithContext(ctx).Preload("User").
		Where("community_id = ?", communityID).
		Offset((page - 1) * limit).Limit(limit).
		Find(&members).Error
	return members, total, err
}

func (r *Repository) GetMember(ctx context.Context, communityID, userID uuid.UUID) (*models.CommunityMember, error) {
	var m models.CommunityMember
	err := database.DB.WithContext(ctx).
		Where("community_id = ? AND user_id = ?", communityID, userID).
		First(&m).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &m, err
}

func (r *Repository) Join(ctx context.Context, communityID, userID uuid.UUID, role string) error {
	m := &models.CommunityMember{CommunityID: communityID, UserID: userID, Role: role}
	if err := database.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("already a member")
	}
	database.DB.WithContext(ctx).Model(&models.Community{}).
		Where("id = ?", communityID).UpdateColumn("member_count", gorm.Expr("member_count + 1"))
	return nil
}

func (r *Repository) Leave(ctx context.Context, communityID, userID uuid.UUID) error {
	result := database.DB.WithContext(ctx).
		Where("community_id = ? AND user_id = ?", communityID, userID).
		Delete(&models.CommunityMember{})
	if result.RowsAffected == 0 {
		return nil
	}
	database.DB.WithContext(ctx).Model(&models.Community{}).
		Where("id = ?", communityID).UpdateColumn("member_count", gorm.Expr("GREATEST(member_count - 1, 0)"))
	return nil
}

func (r *Repository) UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role string) error {
	return database.DB.WithContext(ctx).Model(&models.CommunityMember{}).
		Where("community_id = ? AND user_id = ?", communityID, userID).
		Update("role", role).Error
}

func (r *Repository) GetPosts(ctx context.Context, communityID uuid.UUID, sort string, page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	database.DB.WithContext(ctx).Model(&models.Post{}).Where("community_id = ?", communityID).Count(&total)

	order := "vrs_score DESC"
	switch sort {
	case "latest":
		order = "created_at DESC"
	case "top":
		order = "(like_count + comment_count + share_count) DESC"
	}

	err := database.DB.WithContext(ctx).
		Preload("User").Preload("Tags").
		Where("community_id = ?", communityID).
		Order(order).
		Offset((page - 1) * limit).Limit(limit).
		Find(&posts).Error
	return posts, total, err
}

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
