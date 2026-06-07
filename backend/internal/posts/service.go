package posts

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/internal/engagements"
	"github.com/Nishal77/resona/backend/internal/language"
	"github.com/Nishal77/resona/backend/internal/notifications"
	"github.com/Nishal77/resona/backend/internal/users"
	vrssvc "github.com/Nishal77/resona/backend/internal/vrs"
	"github.com/Nishal77/resona/backend/pkg/models"
)

type Service struct {
	repo          *Repository
	langSvc       *language.Service
	vrsSvc        *vrssvc.Service
	engSvc        *engagements.Service
	notifSvc      *notifications.Service
	usersRepo     *users.Repository
}

func NewService(
	repo *Repository,
	langSvc *language.Service,
	vrsSvc *vrssvc.Service,
	engSvc *engagements.Service,
	notifSvc *notifications.Service,
	usersRepo *users.Repository,
) *Service {
	return &Service{
		repo:      repo,
		langSvc:   langSvc,
		vrsSvc:    vrsSvc,
		engSvc:    engSvc,
		notifSvc:  notifSvc,
		usersRepo: usersRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req *CreatePostRequest) (*models.Post, error) {
	post := &models.Post{
		UserID:    userID,
		ContentText: req.ContentText,
		MediaURL:  req.MediaURL,
		MediaType: req.MediaType,
	}

	if req.CommunityID != nil {
		cid, err := uuid.Parse(*req.CommunityID)
		if err != nil {
			return nil, fmt.Errorf("invalid community_id")
		}
		post.CommunityID = &cid
	}

	// Language detection
	if req.ContentText != nil && *req.ContentText != "" {
		detected, err := s.langSvc.Detect(*req.ContentText)
		if err == nil {
			lang := language.LanguageCodeToName(detected.Language)
			post.DetectedLanguage = &lang
			post.LanguageConfidence = &detected.Confidence
			post.LanguageLocalityScore = detected.LocalityScore
		}
	}

	// Manual language override (when confidence < 0.70)
	if req.ManualLanguage != nil {
		post.DetectedLanguage = req.ManualLanguage
		score := languageToScore(*req.ManualLanguage)
		post.LanguageLocalityScore = score
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	// Initial VRS calculation
	if _, err := s.vrsSvc.Calculate(ctx, post.ID); err != nil {
		_ = err // non-fatal
	}

	return s.repo.FindByID(ctx, post.ID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}
	return post, nil
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *Service) GetFeed(ctx context.Context, userID *uuid.UUID, language string, page, limit int) ([]models.Post, int64, error) {
	if userID == nil {
		return s.repo.GetFeed(ctx, language, page, limit)
	}

	followed := []uuid.UUID{}
	communities := []uuid.UUID{}

	rows, _, _ := s.usersRepo.GetFollowing(ctx, *userID, 1, 1000)
	for _, u := range rows {
		followed = append(followed, u.ID)
	}
	communities, _ = s.usersRepo.GetJoinedCommunityIDs(ctx, *userID)

	if len(followed) == 0 && len(communities) == 0 {
		return s.repo.GetFeed(ctx, language, page, limit)
	}

	return s.repo.GetPersonalizedFeed(ctx, *userID, followed, communities, language, page, limit)
}

func (s *Service) Like(ctx context.Context, postID, userID uuid.UUID) error {
	if err := s.engSvc.Record(ctx, postID, userID, "like"); err != nil {
		return err
	}
	s.repo.IncrementCount(ctx, postID, "like_count", 1)
	s.vrsSvc.Calculate(ctx, postID)
	// notification
	post, _ := s.repo.FindByID(ctx, postID)
	if post != nil && post.UserID != userID {
		s.notifSvc.Create(ctx, post.UserID, userID, "like", &postID, nil)
	}
	return nil
}

func (s *Service) Unlike(ctx context.Context, postID, userID uuid.UUID) error {
	if err := s.engSvc.Delete(ctx, postID, userID, "like"); err != nil {
		return err
	}
	s.repo.DecrementCount(ctx, postID, "like_count")
	s.vrsSvc.Calculate(ctx, postID)
	return nil
}

func (s *Service) Share(ctx context.Context, postID, userID uuid.UUID) error {
	if err := s.engSvc.Record(ctx, postID, userID, "share"); err != nil {
		return err
	}
	s.repo.IncrementCount(ctx, postID, "share_count", 1)
	s.vrsSvc.Calculate(ctx, postID)
	return nil
}

func (s *Service) View(ctx context.Context, postID, userID uuid.UUID) error {
	if err := s.engSvc.Record(ctx, postID, userID, "view"); err != nil {
		return nil // idempotent — duplicate view is not an error to surface
	}
	s.repo.IncrementCount(ctx, postID, "view_count", 1)
	return nil
}

func (s *Service) Save(ctx context.Context, postID, userID uuid.UUID) error {
	if err := s.engSvc.Record(ctx, postID, userID, "save"); err != nil {
		return err
	}
	s.repo.IncrementCount(ctx, postID, "save_count", 1)
	return nil
}

func (s *Service) Unsave(ctx context.Context, postID, userID uuid.UUID) error {
	if err := s.engSvc.Delete(ctx, postID, userID, "save"); err != nil {
		return err
	}
	s.repo.DecrementCount(ctx, postID, "save_count")
	return nil
}

func languageToScore(lang string) float64 {
	scores := map[string]float64{
		"kannada": 1.0, "tamil": 1.0, "telugu": 1.0,
		"malayalam": 1.0, "hindi": 1.0, "english": 0.3,
	}
	if s, ok := scores[lang]; ok {
		return s
	}
	return 0.3
}
