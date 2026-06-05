package explore

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Nishal77/resona/backend/internal/middleware"
	"github.com/Nishal77/resona/backend/pkg/database"
	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/gorm"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Register(r *gin.RouterGroup) {
	e := r.Group("/explore")
	e.GET("/trending", middleware.OptionalAuth(), h.trending)
	e.GET("/tags", h.tags)
	e.GET("/tags/:tagName/posts", h.tagPosts)
	e.GET("/communities", h.communities)
	e.GET("/search", h.search)
}

func (h *Handler) trending(c *gin.Context) {
	lang := c.Query("language")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var posts []models.Post
	q := database.DB.WithContext(c.Request.Context()).
		Preload("User").Preload("Community").Preload("Tags")
	if lang != "" {
		q = q.Where("detected_language = ?", lang)
	}
	q.Order("vrs_score DESC").Limit(limit).Find(&posts)
	middleware.OK(c, posts, "")
}

func (h *Handler) tags(c *gin.Context) {
	lang := c.Query("language")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var tags []models.Tag
	q := database.DB.WithContext(c.Request.Context())
	if lang != "" {
		q = q.Where("language = ?", lang)
	}
	q.Order("trending_score DESC").Limit(limit).Find(&tags)
	middleware.OK(c, tags, "")
}

func (h *Handler) tagPosts(c *gin.Context) {
	tagName := c.Param("tagName")
	page, limit := paginate(c)
	lang := c.Query("language")

	var tag models.Tag
	q := database.DB.WithContext(c.Request.Context()).Where("name = ?", tagName)
	if lang != "" {
		q = q.Where("language = ?", lang)
	}
	if err := q.First(&tag).Error; err != nil {
		middleware.NotFound(c, "tag not found")
		return
	}

	var posts []models.Post
	var total int64
	database.DB.WithContext(c.Request.Context()).
		Model(&models.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tag.ID).
		Count(&total)

	database.DB.WithContext(c.Request.Context()).
		Preload("User").Preload("Community").Preload("Tags").
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tag.ID).
		Order("vrs_score DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&posts)

	middleware.Paginated(c, posts, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
}

func (h *Handler) communities(c *gin.Context) {
	lang := c.Query("language")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "6"))

	var communities []models.Community
	q := database.DB.WithContext(c.Request.Context())
	if lang != "" {
		q = q.Where("primary_language = ?", lang)
	}
	q.Order("member_count DESC").Limit(limit).Find(&communities)
	middleware.OK(c, communities, "")
}

func (h *Handler) search(c *gin.Context) {
	query := c.Query("q")
	lang := c.Query("language")
	page, limit := paginate(c)

	if query == "" {
		middleware.BadRequest(c, "q is required")
		return
	}

	var posts []models.Post
	var total int64
	q := database.DB.WithContext(c.Request.Context()).Model(&models.Post{}).
		Where("content_text ILIKE ?", "%"+query+"%")
	if lang != "" {
		q = q.Where("detected_language = ?", lang)
	}
	q.Count(&total)
	q.Preload("User").Preload("Community").Preload("Tags").
		Order("vrs_score DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&posts)

	middleware.Paginated(c, posts, models.PaginationMeta{
		Page: page, Limit: limit, Total: total,
		HasMore: int64((page-1)*limit+limit) < total,
	})
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

// suppress unused import
var _ = gorm.ErrRecordNotFound
