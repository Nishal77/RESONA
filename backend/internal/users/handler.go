package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/internal/middleware"
	"github.com/Nishal77/resona/backend/pkg/models"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.GET("/me", middleware.AuthRequired(), h.getMe)
	users.PUT("/me", middleware.AuthRequired(), h.updateMe)
	users.GET("/me/communities", middleware.AuthRequired(), h.getMyCommunities)
	users.GET("/me/feed", middleware.AuthRequired(), h.getFeed)
	users.GET("/:username", middleware.OptionalAuth(), h.getByUsername)

	// Routes that need UUID — grouped under /by-id/:id to avoid wildcard conflict
	byID := r.Group("/users/by-id/:id")
	byID.GET("/posts", middleware.OptionalAuth(), h.getPosts)
	byID.GET("/followers", middleware.OptionalAuth(), h.getFollowers)
	byID.GET("/following", middleware.OptionalAuth(), h.getFollowing)
	byID.POST("/follow", middleware.AuthRequired(), h.follow)
	byID.DELETE("/follow", middleware.AuthRequired(), h.unfollow)
}

func (h *Handler) getMe(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	user, err := h.svc.GetMe(c.Request.Context(), userID)
	if err != nil {
		middleware.NotFound(c, err.Error())
		return
	}
	middleware.OK(c, user, "")
}

func (h *Handler) updateMe(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}
	user, err := h.svc.UpdateMe(c.Request.Context(), userID, &req)
	if err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}
	middleware.OK(c, user, "profile updated")
}

func (h *Handler) getByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.svc.GetByUsername(c.Request.Context(), username)
	if err != nil {
		middleware.NotFound(c, "user not found")
		return
	}
	middleware.OK(c, user, "")
}

func (h *Handler) getPosts(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid user id")
		return
	}
	page, limit := paginate(c)
	posts, total, err := h.svc.GetPosts(c.Request.Context(), id, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get posts")
		return
	}
	middleware.Paginated(c, posts, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) getFollowers(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid user id")
		return
	}
	page, limit := paginate(c)
	users, total, err := h.svc.GetFollowers(c.Request.Context(), id, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get followers")
		return
	}
	middleware.Paginated(c, users, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) getFollowing(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid user id")
		return
	}
	page, limit := paginate(c)
	users, total, err := h.svc.GetFollowing(c.Request.Context(), id, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get following")
		return
	}
	middleware.Paginated(c, users, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) follow(c *gin.Context) {
	followerID, _ := middleware.GetUserID(c)
	followingID, err := parseUUID(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid user id")
		return
	}
	if err := h.svc.Follow(c.Request.Context(), followerID, followingID); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}
	middleware.OK(c, nil, "followed")
}

func (h *Handler) unfollow(c *gin.Context) {
	followerID, _ := middleware.GetUserID(c)
	followingID, err := parseUUID(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid user id")
		return
	}
	if err := h.svc.Unfollow(c.Request.Context(), followerID, followingID); err != nil {
		middleware.InternalError(c, err.Error())
		return
	}
	middleware.OK(c, nil, "unfollowed")
}

func (h *Handler) getMyCommunities(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	ids, err := h.svc.GetJoinedCommunities(c.Request.Context(), userID)
	if err != nil {
		middleware.InternalError(c, "failed to get communities")
		return
	}
	middleware.OK(c, ids, "")
}

func (h *Handler) getFeed(c *gin.Context) {
	// Feed is handled by posts service — this is a redirect alias
	c.Redirect(302, "/api/v1/posts?"+c.Request.URL.RawQuery)
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func paginate(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
