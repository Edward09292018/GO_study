package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func HandleError(c *gin.Context, err error, statusCode int) {
	log.Printf("Error: %v", err)
	c.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
}

func LogRequest(c *gin.Context) {
	log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SendSuccess(c *gin.Context, data interface{}, message string, statusCode int) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

var Logger *zap.Logger

func LogInfo(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(message, fields...)
	}
}
