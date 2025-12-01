package utils

import "github.com/gin-gonic/gin"

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func RespondOK(c *gin.Context, data map[string]interface{}) {
	c.JSON(200, data)
}
