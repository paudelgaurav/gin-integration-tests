package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudelgaurav/gin-integration-tests/pkg/framework"
)

func HandleValidationError(c *gin.Context, logger framework.Logger, obj any, err error) {
	logger.Error(err)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": formatValidationError(err, obj),
	})
}

func HandleErrorWithStatus(c *gin.Context, logger framework.Logger, statusCode int, err error) {
	logger.Error(err)
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
