package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository { return &Repository{} }

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	return database.DB.WithContext(ctx).Create(user).Error
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := database.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
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

func (r *Repository) FindByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := database.DB.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := database.DB.WithContext(ctx).First(&user, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	token := &models.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return database.DB.WithContext(ctx).Create(token).Error
}

func (r *Repository) FindRefreshToken(ctx context.Context, userID uuid.UUID) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := database.DB.WithContext(ctx).
		Where("user_id = ? AND expires_at > NOW()", userID).
		Order("created_at DESC").
		First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &token, err
}

func (r *Repository) DeleteRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return database.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).Error
}

func (r *Repository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check username: %w", err)
	}
	return count > 0, nil
}
