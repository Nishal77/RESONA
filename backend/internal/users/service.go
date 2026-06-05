package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) GetMe(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Service) UpdateMe(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.FullName != nil {
		user.FullName = req.FullName
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.State != nil {
		user.State = req.State
	}
	if req.City != nil {
		user.City = req.City
	}
	if req.PrimaryLanguage != nil {
		user.PrimaryLanguage = *req.PrimaryLanguage
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return user, nil
}

func (s *Service) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}
	following, err := s.repo.FindByID(ctx, followingID)
	if err != nil || following == nil {
		return fmt.Errorf("user not found")
	}
	return s.repo.Follow(ctx, followerID, followingID)
}

func (s *Service) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.repo.Unfollow(ctx, followerID, followingID)
}

func (s *Service) GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.User, int64, error) {
	return s.repo.GetFollowers(ctx, userID, page, limit)
}

func (s *Service) GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.User, int64, error) {
	return s.repo.GetFollowing(ctx, userID, page, limit)
}

func (s *Service) GetPosts(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Post, int64, error) {
	return s.repo.GetPostsByUser(ctx, userID, page, limit)
}

func (s *Service) GetJoinedCommunities(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.GetJoinedCommunityIDs(ctx, userID)
}
