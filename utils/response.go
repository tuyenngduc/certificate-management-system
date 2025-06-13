package utils

import (
	"github.com/gin-gonic/gin"
)

// Success response helper
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"data": data,
	})
}

// Error response helper
func Error(c *gin.Context, statusCode int, message string, err error) {
	resp := gin.H{
		"error": message,
	}
	if err != nil {
		resp["details"] = err.Error()
	}
	c.JSON(statusCode, resp)
}
