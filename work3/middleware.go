package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const RequestIDKey = "X-Request-ID"

// JWT 密钥建议从配置中读取
var jwtSecret = os.Getenv("JWT_SECRET") // 可替换为 os.Getenv("JWT_SECRET")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("cookie")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "missing authorization header"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token"})
			return
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token claims"})
			return
		}

		c.Set("username", (*claims)["username"])
		c.Set("userid", (*claims)["id"])
		c.Next()
	}
}

func AuthorOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("userid")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
			return
		}

		postIDStr := c.Param("postId")
		if postIDStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "missing postId in URL path"})
			return
		}

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "invalid postId format"})
			return
		}

		var post Post
		if err := db.Debug().Where("id = ?", uint(postID)).First(&post).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.WithField("request_id", c.MustGet("X-Request-ID")).Warn("Post not found")
				c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: "post not found"})
			} else {
				logrus.WithFields(logrus.Fields{
					"request_id": c.MustGet("X-Request-ID"),
					"error":      err,
				}).Error("Database error")
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "database error"})
			}
			return
		}

		currentUserID, ok := userIDInterface.(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid user ID type"})
			return
		}

		if uint(currentUserID) != post.UserID {
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: "you are not the author of this post"})
			return
		}

		c.Set("post", &post)
		c.Next()
	}
}

func CustomErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("X-Request-ID")
				logrus.WithFields(logrus.Fields{
					"request_id": requestID,
					"error":      err,
				}).Error("Panic occurred")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get(RequestIDKey)
		userID, _ := c.Get("userid")

		fields := logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}

		if requestID != nil {
			fields[RequestIDKey] = requestID
		}
		if userID != nil {
			fields["user_id"] = userID
		}

		// 记录开始日志
		logrus.WithFields(fields).Info("Request started")

		c.Next()

		// 获取响应状态码
		status := c.Writer.Status()
		fields["status"] = status

		// 记录结束日志
		if status >= 500 {
			logrus.WithFields(fields).Error("Request completed with error")
		} else if status >= 400 {
			logrus.WithFields(fields).Warn("Request completed with warning")
		} else {
			logrus.WithFields(fields).Info("Request completed")
		}

		// 构造日志对象
		log := &RequestLog{
			RequestID: requestID.(string),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			IP:        c.ClientIP(),
			Status:    status,
		}

		if userID != nil {
			log.UserID = userID.(string)
		}

		// 写入数据库
		if err := db.Create(log).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"error":      err,
				"request_id": requestID,
			}).Error("Failed to save request log to database")
		}
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDKey, requestID)
		c.Next()
	}
}
