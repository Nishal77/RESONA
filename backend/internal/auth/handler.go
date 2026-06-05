package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Nishal77/resona/backend/internal/middleware"
	"github.com/Nishal77/resona/backend/pkg/config"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/register", h.register)
	auth.POST("/login", h.login)
	auth.POST("/google", h.googleAuth)
	auth.POST("/refresh", h.refresh)
	auth.POST("/logout", middleware.AuthRequired(), h.logout)
}

func (h *Handler) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	setRefreshCookie(c, resp.RefreshToken)
	middleware.Created(c, gin.H{"access_token": resp.AccessToken, "user": resp.User}, "registered successfully")
}

func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		middleware.Unauthorized(c, err.Error())
		return
	}

	setRefreshCookie(c, resp.RefreshToken)
	middleware.OK(c, gin.H{"access_token": resp.AccessToken, "user": resp.User}, "login successful")
}

func (h *Handler) googleAuth(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	resp, err := h.svc.GoogleAuth(c.Request.Context(), req.GoogleToken)
	if err != nil {
		middleware.Unauthorized(c, err.Error())
		return
	}

	setRefreshCookie(c, resp.RefreshToken)
	middleware.OK(c, gin.H{"access_token": resp.AccessToken, "user": resp.User}, "google auth successful")
}

func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		middleware.Unauthorized(c, "refresh token missing")
		return
	}

	accessToken, err := h.svc.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		middleware.Unauthorized(c, err.Error())
		return
	}

	middleware.OK(c, gin.H{"access_token": accessToken}, "token refreshed")
}

func (h *Handler) logout(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		middleware.Unauthorized(c, "unauthorized")
		return
	}

	if err := h.svc.Logout(c.Request.Context(), userID); err != nil {
		middleware.InternalError(c, "logout failed")
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	middleware.OK(c, nil, "logged out")
}

func setRefreshCookie(c *gin.Context, token string) {
	maxAge := int(config.App.JWTRefreshExpiresIn / time.Second)
	secure := config.App.AppEnv == "production"
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", token, maxAge, "/", "", secure, true)
}
