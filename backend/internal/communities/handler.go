package communities

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/internal/middleware"
	"github.com/Nishal77/resona/backend/pkg/models"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Register(r *gin.RouterGroup) {
	c := r.Group("/communities")
	c.GET("", middleware.OptionalAuth(), h.list)
	c.POST("", middleware.AuthRequired(), h.create)
	c.GET("/:slug", middleware.OptionalAuth(), h.getBySlug)

	byID := c.Group("/by-id/:id")
	byID.PUT("", middleware.AuthRequired(), h.update)
	byID.GET("/posts", middleware.OptionalAuth(), h.getPosts)
	byID.GET("/members", middleware.OptionalAuth(), h.getMembers)
	byID.GET("/snap-of-week", middleware.OptionalAuth(), h.snapOfWeek)
	byID.POST("/join", middleware.AuthRequired(), h.join)
	byID.DELETE("/join", middleware.AuthRequired(), h.leave)
	byID.PUT("/members/:userId/role", middleware.AuthRequired(), h.updateRole)
}

func (h *Handler) list(c *gin.Context) {
	page, limit := paginate(c)
	lang := c.Query("language")
	communities, total, err := h.repo.List(c.Request.Context(), lang, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get communities")
		return
	}
	middleware.Paginated(c, communities, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var req CreateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	slug := GenerateSlug(req.Name)
	community := &models.Community{
		Name:               req.Name,
		Slug:               slug,
		Description:        req.Description,
		PrimaryLanguage:    req.PrimaryLanguage,
		SecondaryLanguages: req.SecondaryLanguages,
		AvatarURL:          req.AvatarURL,
		BannerURL:          req.BannerURL,
		CreatedBy:          &userID,
	}

	if err := h.repo.Create(c.Request.Context(), community); err != nil {
		middleware.BadRequest(c, "community name/slug already taken")
		return
	}

	// Creator becomes admin
	h.repo.Join(c.Request.Context(), community.ID, userID, "admin")

	middleware.Created(c, community, "community created")
}

func (h *Handler) getBySlug(c *gin.Context) {
	slug := c.Param("slug")
	community, err := h.repo.FindBySlug(c.Request.Context(), slug)
	if err != nil || community == nil {
		middleware.NotFound(c, "community not found")
		return
	}
	middleware.OK(c, community, "")
}

func (h *Handler) update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid id")
		return
	}

	member, err := h.repo.GetMember(c.Request.Context(), id, userID)
	if err != nil || member == nil || member.Role != "admin" {
		middleware.Forbidden(c, "must be community admin")
		return
	}

	var req UpdateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	community, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil || community == nil {
		middleware.NotFound(c, "community not found")
		return
	}

	if req.Name != nil {
		community.Name = *req.Name
	}
	if req.Description != nil {
		community.Description = req.Description
	}
	if req.SecondaryLanguages != nil {
		community.SecondaryLanguages = req.SecondaryLanguages
	}
	if req.AvatarURL != nil {
		community.AvatarURL = req.AvatarURL
	}
	if req.BannerURL != nil {
		community.BannerURL = req.BannerURL
	}

	if err := h.repo.Update(c.Request.Context(), community); err != nil {
		middleware.InternalError(c, "update failed")
		return
	}
	middleware.OK(c, community, "community updated")
}

func (h *Handler) getPosts(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid id")
		return
	}
	page, limit := paginate(c)
	sort := c.DefaultQuery("sort", "vrs")
	posts, total, err := h.repo.GetPosts(c.Request.Context(), id, sort, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get posts")
		return
	}
	middleware.Paginated(c, posts, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) getMembers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid id")
		return
	}
	page, limit := paginate(c)
	members, total, err := h.repo.GetMembers(c.Request.Context(), id, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get members")
		return
	}
	middleware.Paginated(c, members, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) snapOfWeek(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid id")
		return
	}
	community, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil || community == nil {
		middleware.NotFound(c, "community not found")
		return
	}
	middleware.OK(c, gin.H{"snap_of_week_post_id": community.SnapOfWeekPostID}, "")
}

func (h *Handler) join(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid id")
		return
	}
	if err := h.repo.Join(c.Request.Context(), id, userID, "member"); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}
	middleware.OK(c, nil, "joined")
}

func (h *Handler) leave(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid id")
		return
	}
	h.repo.Leave(c.Request.Context(), id, userID)
	middleware.OK(c, nil, "left")
}

func (h *Handler) updateRole(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	communityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid community id")
		return
	}
	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		middleware.BadRequest(c, "invalid user id")
		return
	}

	member, err := h.repo.GetMember(c.Request.Context(), communityID, userID)
	if err != nil || member == nil || member.Role != "admin" {
		middleware.Forbidden(c, "must be community admin")
		return
	}

	var req UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	if err := h.repo.UpdateMemberRole(c.Request.Context(), communityID, targetUserID, req.Role); err != nil {
		middleware.InternalError(c, "failed to update role")
		return
	}
	middleware.OK(c, nil, "role updated")
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
