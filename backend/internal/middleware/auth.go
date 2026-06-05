package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/pkg/config"
)

const UserIDKey = "userID"

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			Unauthorized(c, "missing authorization token")
			c.Abort()
			return
		}

		claims, err := validateAccessToken(token)
		if err != nil {
			Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		id, err := uuid.Parse(claims.UserID)
		if err != nil {
			Unauthorized(c, "invalid token claims")
			c.Abort()
			return
		}

		c.Set(UserIDKey, id)
		c.Next()
	}
}

func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			claims, err := validateAccessToken(token)
			if err == nil {
				if id, err := uuid.Parse(claims.UserID); err == nil {
					c.Set(UserIDKey, id)
				}
			}
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}

func validateAccessToken(tokenStr string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.App.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
