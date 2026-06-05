package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Nishal77/resona/backend/pkg/models"
)

func OK(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Paginated(c *gin.Context, data interface{}, meta models.PaginationMeta) {
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, models.APIResponse{
		Success:    false,
		Error:      err,
		StatusCode: http.StatusBadRequest,
	})
}

func Unauthorized(c *gin.Context, err string) {
	c.JSON(http.StatusUnauthorized, models.APIResponse{
		Success:    false,
		Error:      err,
		StatusCode: http.StatusUnauthorized,
	})
}

func Forbidden(c *gin.Context, err string) {
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success:    false,
		Error:      err,
		StatusCode: http.StatusForbidden,
	})
}

func NotFound(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, models.APIResponse{
		Success:    false,
		Error:      err,
		StatusCode: http.StatusNotFound,
	})
}

func InternalError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, models.APIResponse{
		Success:    false,
		Error:      err,
		StatusCode: http.StatusInternalServerError,
	})
}
