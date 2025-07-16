package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var post Post
	result := db.WithContext(c).Where("id = ? AND deleted_at IS NULL", comment.PostID).Find(&post)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Post not found"})
		return
	}

	userID, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	uid, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid user ID type"})
		return
	}
	comment.UserID = uint(uid)

	result = db.WithContext(c).Preload("User").Preload("Post").Create(&comment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create comment"})
		return
	}

	if comment.Post == nil || comment.User == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to load related data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"title":   comment.Post.Title,
		"content": comment.Content,
		"user":    comment.User.Username,
	}})
}

func ReadComment(c *gin.Context) {
	postIDStr := c.Param("postId")
	if postIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing post ID"})
		return
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid post ID format"})
		return
	}

	var comments []*Comment
	db.WithContext(c).Where("post_id = ?", postID).Find(&comments)

	c.JSON(http.StatusOK, gin.H{"data": comments})
}
