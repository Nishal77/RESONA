package upload

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Nishal77/resona/backend/internal/middleware"
	"github.com/Nishal77/resona/backend/pkg/supabase"
)

type Handler struct {
	storage *supabase.StorageClient
}

func NewHandler(storage *supabase.StorageClient) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	upload := r.Group("/upload", middleware.AuthRequired())
	upload.POST("/image", h.uploadImage)
	upload.POST("/video", h.uploadVideo)
}

func (h *Handler) uploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		middleware.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	if header.Size > 10<<20 { // 10MB
		middleware.BadRequest(c, "image must be under 10MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		middleware.BadRequest(c, "allowed formats: jpg, png, webp")
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		middleware.InternalError(c, "failed to read file")
		return
	}

	url, err := h.storage.Upload(content, header.Filename, fmt.Sprintf("image/%s", ext[1:]))
	if err != nil {
		middleware.InternalError(c, "upload failed")
		return
	}

	middleware.OK(c, gin.H{"url": url}, "uploaded")
}

func (h *Handler) uploadVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		middleware.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	if header.Size > 50<<20 { // 50MB
		middleware.BadRequest(c, "video must be under 50MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".mp4": true, ".webm": true}
	if !allowed[ext] {
		middleware.BadRequest(c, "allowed formats: mp4, webm")
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		middleware.InternalError(c, "failed to read file")
		return
	}

	url, err := h.storage.Upload(content, header.Filename, fmt.Sprintf("video/%s", ext[1:]))
	if err != nil {
		middleware.InternalError(c, "upload failed")
		return
	}

	middleware.OK(c, gin.H{"url": url}, "uploaded")
}
