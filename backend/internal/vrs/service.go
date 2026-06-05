package vrs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	rdb "github.com/Nishal77/resona/backend/pkg/redis"
	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service { return &Service{} }

// Calculate computes and persists VRS for a single post.
// VRS = (EngagementRate × LocalityScore × ShareVelocity) / (HoursSincePosted + 1)
func (s *Service) Calculate(ctx context.Context, postID uuid.UUID) (float64, error) {
	var post models.Post
	if err := database.DB.WithContext(ctx).First(&post, "id = ?", postID).Error; err != nil {
		return 0, fmt.Errorf("fetch post: %w", err)
	}

	var engagementRate float64
	if post.ViewCount > 0 {
		engagementRate = float64(post.LikeCount+post.CommentCount+post.ShareCount) / float64(post.ViewCount)
	}

	// Share velocity: shares in first 2 hours after posting
	velocityWindow := post.CreatedAt.Add(2 * time.Hour)
	var shareVelocity int64
	database.DB.WithContext(ctx).Model(&models.Engagement{}).
		Where("post_id = ? AND type = 'share' AND created_at <= ?", postID, velocityWindow).
		Count(&shareVelocity)

	hoursSince := time.Since(post.CreatedAt).Hours()
	vrs := (engagementRate * post.LanguageLocalityScore * float64(shareVelocity)) / (hoursSince + 1)

	if err := database.DB.WithContext(ctx).Model(&post).Update("vrs_score", vrs).Error; err != nil {
		return 0, fmt.Errorf("update vrs_score: %w", err)
	}

	return vrs, nil
}

// RecalculateAll runs on the cron schedule — updates every post and user vrs_total.
func (s *Service) RecalculateAll(ctx context.Context) error {
	var postIDs []uuid.UUID
	if err := database.DB.WithContext(ctx).Model(&models.Post{}).Pluck("id", &postIDs).Error; err != nil {
		return fmt.Errorf("fetch post ids: %w", err)
	}

	batchSize := 100
	for i := 0; i < len(postIDs); i += batchSize {
		end := i + batchSize
		if end > len(postIDs) {
			end = len(postIDs)
		}
		batch := postIDs[i:end]
		for _, id := range batch {
			if _, err := s.Calculate(ctx, id); err != nil {
				// log and continue — one bad post shouldn't halt all recalculations
				_ = err
			}
		}
	}

	// Update vrs_total on each user
	if err := database.DB.WithContext(ctx).Exec(`
		UPDATE users u
		SET vrs_total = COALESCE((
			SELECT SUM(p.vrs_score) FROM posts p WHERE p.user_id = u.id
		), 0)
	`).Error; err != nil {
		return fmt.Errorf("update user vrs_total: %w", err)
	}

	// Invalidate feed cache
	if err := rdb.DelPattern(ctx, "feed:*"); err != nil {
		_ = err // non-fatal
	}
	if err := rdb.DelPattern(ctx, "explore:*"); err != nil {
		_ = err
	}

	return nil
}

// UpdateTrendingTags updates tag trending_score based on usage in last 24h.
func (s *Service) UpdateTrendingTags(ctx context.Context) error {
	return database.DB.WithContext(ctx).Exec(`
		UPDATE tags t
		SET trending_score = COALESCE((
			SELECT COUNT(pt.post_id)
			FROM post_tags pt
			JOIN posts p ON p.id = pt.post_id
			WHERE pt.tag_id = t.id AND p.created_at >= NOW() - INTERVAL '24 hours'
		), 0)
	`).Error
}

// SetSnapOfWeek picks highest VRS post per community from the last 7 days.
// Only runs on Monday at 00:00.
func (s *Service) SetSnapOfWeek(ctx context.Context) error {
	return database.DB.WithContext(ctx).Exec(`
		UPDATE communities c
		SET snap_of_week_post_id = (
			SELECT p.id
			FROM posts p
			WHERE p.community_id = c.id
			  AND p.created_at >= NOW() - INTERVAL '7 days'
			ORDER BY p.vrs_score DESC
			LIMIT 1
		)
	`).Error
}

// CreateTrendingNotifications fires notifications for posts crossing the threshold.
func (s *Service) CreateTrendingNotifications(ctx context.Context, threshold float64) error {
	var posts []models.Post
	if err := database.DB.WithContext(ctx).
		Where("vrs_score >= ? AND community_id IS NOT NULL", threshold).
		Find(&posts).Error; err != nil {
		return err
	}

	for _, p := range posts {
		msg := "Your post is trending!"
		existing := &models.Notification{}
		err := database.DB.WithContext(ctx).
			Where("post_id = ? AND type = 'trending'", p.ID).
			First(existing).Error
		if err == gorm.ErrRecordNotFound {
			database.DB.WithContext(ctx).Create(&models.Notification{
				UserID:  p.UserID,
				Type:    "trending",
				PostID:  &p.ID,
				Message: &msg,
			})
		}
	}
	return nil
}
