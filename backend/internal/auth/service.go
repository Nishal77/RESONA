package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/config"
	"github.com/Nishal77/resona/backend/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword SHA-256 prehashes the password before bcrypt.
// bcrypt silently truncates input at 72 bytes — prehashing avoids that entirely.
func hashPassword(plain string) []byte {
	sum := sha256.Sum256([]byte(plain))
	// hex-encode so result is 64 printable ASCII bytes — safe for bcrypt
	encoded := hex.EncodeToString(sum[:])
	return []byte(encoded)
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	exists, err := s.repo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("username already taken")
	}

	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword(hashPassword(req.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	hashStr := string(hash)
	fullName := req.FullName
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &hashStr,
		FullName:     &fullName,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.PasswordHash == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), hashPassword(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *Service) GoogleAuth(ctx context.Context, googleToken string) (*AuthResponse, error) {
	googleUser, err := verifyGoogleToken(googleToken)
	if err != nil {
		return nil, fmt.Errorf("invalid google token: %w", err)
	}

	user, err := s.repo.FindByGoogleID(ctx, googleUser.Sub)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Try to find by email first
		user, err = s.repo.FindByEmail(ctx, googleUser.Email)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		// New user via Google
		username := generateUsername(googleUser.Email)
		user = &models.User{
			Username:    username,
			Email:       googleUser.Email,
			GoogleID:    &googleUser.Sub,
			FullName:    &googleUser.Name,
			AvatarURL:   &googleUser.Picture,
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("create google user: %w", err)
		}
	} else if user.GoogleID == nil {
		// Link Google to existing email account
		user.GoogleID = &googleUser.Sub
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.App.JWTRefreshSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("invalid token claims")
	}

	stored, err := s.repo.FindRefreshToken(ctx, userID)
	if err != nil || stored == nil {
		return "", fmt.Errorf("refresh token not found or expired")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored.TokenHash), hashPassword(refreshToken)); err != nil {
		return "", fmt.Errorf("refresh token mismatch")
	}

	// Rotate: delete old, issue new is handled by caller
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return "", fmt.Errorf("user not found")
	}

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *Service) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteRefreshTokens(ctx, userID)
}

func (s *Service) buildAuthResponse(ctx context.Context, user *models.User) (*AuthResponse, error) {
	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(hashPassword(refreshToken), 10)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(config.App.JWTRefreshExpiresIn)
	if err := s.repo.SaveRefreshToken(ctx, user.ID, string(hash), expiresAt); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"-"` // sent via httpOnly cookie, never in JSON
	User         *models.User `json:"user"`
}

func generateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwtClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.App.JWTExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.App.JWTSecret))
}

func generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwtClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.App.JWTRefreshExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.App.JWTRefreshSecret))
}

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// verifyGoogleToken accepts either:
// - an OAuth2 access_token (from implicit/popup flow) → calls userinfo endpoint
// - an id_token (credential from One Tap) → calls tokeninfo endpoint
// Frontend sends access_token via useGoogleLogin implicit flow.
func verifyGoogleToken(token string) (*googleUserInfo, error) {
	// Try userinfo first (access_token from implicit flow)
	req, _ := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var info googleUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return nil, fmt.Errorf("decode userinfo: %w", err)
		}
		if info.Sub == "" || info.Email == "" {
			return nil, fmt.Errorf("google userinfo: missing sub or email")
		}
		return &info, nil
	}

	// Fallback: try tokeninfo (id_token from One Tap credential)
	resp2, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + token)
	if err != nil {
		return nil, fmt.Errorf("google tokeninfo: %w", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google auth failed: status %d", resp2.StatusCode)
	}
	var info googleUserInfo
	if err := json.NewDecoder(resp2.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode tokeninfo: %w", err)
	}
	return &info, nil
}

func generateUsername(email string) string {
	parts := strings.Split(email, "@")
	base := parts[0]
	base = strings.ReplaceAll(base, ".", "_")
	if len(base) > 40 {
		base = base[:40]
	}
	return base + "_" + uuid.New().String()[:4]
}
