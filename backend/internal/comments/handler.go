package comments

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Nishal77/resona/backend/internal/middleware"
	"github.com/Nishal77/resona/backend/internal/notifications"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/gorm"
)

type Handler struct {
	repo     *Repository
	notifSvc *notifications.Service
}

func NewHandler(repo *Repository, notifSvc *notifications.Service) *Handler {
	return &Handler{repo: repo, notifSvc: notifSvc}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	posts := r.Group("/posts")
	posts.GET("/:id/comments", middleware.OptionalAuth(), h.list)
	posts.POST("/:id/comments", middleware.AuthRequired(), h.create)
	posts.DELETE("/:id/comments/:commentId", middleware.AuthRequired(), h.delete)
	posts.POST("/:id/comments/:commentId/like", middleware.AuthRequired(), h.like)
}

func (h *Handler) list(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid post id")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	comments, total, err := h.repo.GetByPost(c.Request.Context(), postID, page, limit)
	if err != nil {
		middleware.InternalError(c, "failed to get comments")
		return
	}
	middleware.Paginated(c, comments, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.BadRequest(c, "invalid post id")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequest(c, err.Error())
		return
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	if req.ParentCommentID != nil {
		pid, err := uuid.Parse(*req.ParentCommentID)
		if err != nil {
			middleware.BadRequest(c, "invalid parent_comment_id")
			return
		}
		comment.ParentCommentID = &pid
	}

	if err := h.repo.Create(c.Request.Context(), comment); err != nil {
		middleware.InternalError(c, "failed to create comment")
		return
	}

	// increment post.comment_count
	database.DB.WithContext(c.Request.Context()).Model(&models.Post{}).
		Where("id = ?", postID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	// notification to post owner
	h.notifSvc.Create(c.Request.Context(), postID, userID, "comment", &postID, nil)

	middleware.Created(c, comment, "comment added")
}

func (h *Handler) delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	commentID, err := uuid.Parse(c.Param("commentId"))
	if err != nil {
		middleware.BadRequest(c, "invalid comment id")
		return
	}
	postID, _ := uuid.Parse(c.Param("id"))

	if err := h.repo.Delete(c.Request.Context(), commentID, userID); err != nil {
		middleware.Forbidden(c, err.Error())
		return
	}

	database.DB.WithContext(c.Request.Context()).Model(&models.Post{}).
		Where("id = ?", postID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))

	middleware.OK(c, nil, "comment deleted")
}

func (h *Handler) like(c *gin.Context) {
	commentID, err := uuid.Parse(c.Param("commentId"))
	if err != nil {
		middleware.BadRequest(c, "invalid comment id")
		return
	}
	if err := h.repo.IncrementLike(c.Request.Context(), commentID); err != nil {
		middleware.InternalError(c, "failed to like comment")
		return
	}
	middleware.OK(c, nil, "liked")
}

// notification helper — resolves post owner
func (h *Handler) notifyPostOwner(ctx context.Context, postID, actorID uuid.UUID, notifType string) {
	var post models.Post
	if err := database.DB.WithContext(ctx).Select("user_id").First(&post, "id = ?", postID).Error; err != nil {
		return
	}
	if post.UserID != actorID {
		h.notifSvc.Create(ctx, post.UserID, actorID, notifType, &postID, nil)
	}
}

// suppress unused import guards
var (
	_ = fmt.Sprintf
	_ = gorm.ErrRecordNotFound
)
