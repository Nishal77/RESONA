package notifications

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
	n := r.Group("/notifications", middleware.AuthRequired())
	n.GET("", h.list)
	n.GET("/unread-count", h.unreadCount)
	n.PUT("/:id/read", h.markRead)
	n.PUT("/read-all", h.markAllRead)
}

func (h *Handler) list(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	notifs, total, err := h.svc.List(c.Request.Context(), userID, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get notifications")
		return
	}
	middleware.Paginated(c, notifs, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) unreadCount(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	count, err := h.svc.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		middleware.InternalError(c, "failed to get count")
		return
	}
	middleware.OK(c, gin.H{"count": count}, "")
}

func (h *Handler) markRead(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.MarkRead(c.Request.Context(), id, userID); err != nil {
		middleware.NotFound(c, err.Error())
		return
	}
	middleware.OK(c, nil, "marked as read")
}

func (h *Handler) markAllRead(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	updated, err := h.svc.MarkAllRead(c.Request.Context(), userID)
	if err != nil {
		middleware.InternalError(c, "failed to mark all read")
		return
	}
	middleware.OK(c, gin.H{"updated": updated}, "all marked as read")
}
