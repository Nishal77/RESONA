package posts

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
	posts := r.Group("/posts")
	posts.GET("", middleware.OptionalAuth(), h.getFeed)
	posts.POST("", middleware.AuthRequired(), h.create)
	posts.GET("/:id", middleware.OptionalAuth(), h.getByID)
	posts.DELETE("/:id", middleware.AuthRequired(), h.delete)
	posts.POST("/:id/like", middleware.AuthRequired(), h.like)
	posts.DELETE("/:id/like", middleware.AuthRequired(), h.unlike)
	posts.POST("/:id/share", middleware.AuthRequired(), h.share)
	posts.POST("/:id/view", middleware.OptionalAuth(), h.view)
	posts.POST("/:id/save", middleware.AuthRequired(), h.save)
	posts.DELETE("/:id/save", middleware.AuthRequired(), h.unsave)
}

func (h *Handler) getFeed(c *gin.Context) {
	page, limit := paginate(c)
	lang := c.Query("language")
	userID, ok := middleware.GetUserID(c)
	var uid *uuid.UUID
	if ok {
		uid = &userID
	}
	posts, total, err := h.svc.GetFeed(c.Request.Context(), uid, lang, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get feed")
		return
	}
	middleware.Paginated(c, posts, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}
	if req.ContentText == nil && req.MediaURL == nil {
		middleware.BadRequest(c, "post must have content or media")
		return
	}
	post, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		middleware.InternalError(c, err.Error())
		return
	}
	middleware.Created(c, post, "post created")
}

func (h *Handler) getByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid post id")
		return
	}
	post, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		middleware.NotFound(c, "post not found")
		return
	}
	middleware.OK(c, post, "")
}

func (h *Handler) delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid post id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		middleware.Forbidden(c, err.Error())
		return
	}
	middleware.OK(c, nil, "post deleted")
}

func (h *Handler) like(c *gin.Context) {
	h.engagement(c, func(postID, userID uuid.UUID) error {
		return h.svc.Like(c.Request.Context(), postID, userID)
	})
}

func (h *Handler) unlike(c *gin.Context) {
	h.engagement(c, func(postID, userID uuid.UUID) error {
		return h.svc.Unlike(c.Request.Context(), postID, userID)
	})
}

func (h *Handler) share(c *gin.Context) {
	h.engagement(c, func(postID, userID uuid.UUID) error {
		return h.svc.Share(c.Request.Context(), postID, userID)
	})
}

func (h *Handler) view(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	userID, ok := middleware.GetUserID(c)
	if !ok {
		middleware.OK(c, nil, "")
		return
	}
	h.svc.View(c.Request.Context(), id, userID)
	middleware.OK(c, nil, "")
}

func (h *Handler) save(c *gin.Context) {
	h.engagement(c, func(postID, userID uuid.UUID) error {
		return h.svc.Save(c.Request.Context(), postID, userID)
	})
}

func (h *Handler) unsave(c *gin.Context) {
	h.engagement(c, func(postID, userID uuid.UUID) error {
		return h.svc.Unsave(c.Request.Context(), postID, userID)
	})
}

func (h *Handler) engagement(c *gin.Context, fn func(uuid.UUID, uuid.UUID) error) {
	userID, _ := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid post id")
		return
	}
	if err := fn(id, userID); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}
	middleware.OK(c, nil, "ok")
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
