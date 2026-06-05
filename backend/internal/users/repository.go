package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository { return &Repository{} }

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := database.DB.WithContext(ctx).First(&user, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := database.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) Update(ctx context.Context, user *models.User) error {
	return database.DB.WithContext(ctx).Save(user).Error
}

func (r *Repository) GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	subQuery := database.DB.Model(&models.Follow{}).Select("follower_id").Where("following_id = ?", userID)
	database.DB.WithContext(ctx).Model(&models.User{}).Where("id IN (?)", subQuery).Count(&total)
	err := database.DB.WithContext(ctx).Where("id IN (?)", subQuery).
		Offset((page - 1) * limit).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *Repository) GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	subQuery := database.DB.Model(&models.Follow{}).Select("following_id").Where("follower_id = ?", userID)
	database.DB.WithContext(ctx).Model(&models.User{}).Where("id IN (?)", subQuery).Count(&total)
	err := database.DB.WithContext(ctx).Where("id IN (?)", subQuery).
		Offset((page - 1) * limit).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *Repository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	follow := models.Follow{FollowerID: followerID, FollowingID: followingID}
	if err := database.DB.WithContext(ctx).Create(&follow).Error; err != nil {
		return err
	}
	database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", followerID).
		UpdateColumn("following_count", gorm.Expr("following_count + 1"))
	database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", followingID).
		UpdateColumn("follower_count", gorm.Expr("follower_count + 1"))
	return nil
}

func (r *Repository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	result := database.DB.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.Follow{})
	if result.RowsAffected == 0 {
		return nil
	}
	database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", followerID).
		UpdateColumn("following_count", gorm.Expr("GREATEST(following_count - 1, 0)"))
	database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", followingID).
		UpdateColumn("follower_count", gorm.Expr("GREATEST(follower_count - 1, 0)"))
	return nil
}

func (r *Repository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) GetJoinedCommunityIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := database.DB.WithContext(ctx).Model(&models.CommunityMember{}).
		Where("user_id = ?", userID).Pluck("community_id", &ids).Error
	return ids, err
}

func (r *Repository) GetPostsByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	database.DB.WithContext(ctx).Model(&models.Post{}).Where("user_id = ?", userID).Count(&total)
	err := database.DB.WithContext(ctx).
		Preload("User").Preload("Community").Preload("Tags").
		Where("user_id = ?", userID).
		Order("vrs_score DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&posts).Error
	return posts, total, err
}
