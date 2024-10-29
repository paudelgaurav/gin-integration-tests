package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, gin.H{"data": data})
}

func ValidationError(c *gin.Context, errMessage string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": errMessage})
}

func InternalServerError(c *gin.Context, errMessage string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
}
