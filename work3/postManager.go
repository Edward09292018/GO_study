package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FileCreated(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	// 安全类型断言
	userIDFloat, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid user ID format"})
		return
	}
	post.UserID = uint(userIDFloat)

	result := db.Preload("User").FirstOrCreate(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func FileRead(c *gin.Context) {
	var post Post
	postID := c.Param("postId")
	id, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid postId format"})
		return
	}

	result := db.Select("title", "content").Where("id = ?", id).Find(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to query post"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

// 命名结构体，提高可维护性
type PostListItem struct {
	ID    uint   `json:"postId" gorm:"column:id"`
	Title string `json:"title" gorm:"column:title"`
}

func FileList(c *gin.Context) {
	var posts []PostListItem

	result := db.Table("posts").Select("id, title").Where("deleted_at is null").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch post list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func FileUpdated(c *gin.Context) {
	var postNew Post
	if err := c.ShouldBindJSON(&postNew); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// 从中间件获取 post
	rawPost, exists := c.Get("post")
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Post not found in context"})
		return
	}
	post, ok := rawPost.(*Post)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid post type in context"})
		return
	}

	// 获取用户ID用于权限校验
	userID, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}
	userIDFloat, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	// 权限校验
	if post.UserID != uint(userIDFloat) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "You are not allowed to update this post"})
		return
	}

	post.Title = postNew.Title
	post.Content = postNew.Content

	result := db.Save(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func FileDeleted(c *gin.Context) {
	rawPost, exists := c.Get("post")
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Post not found in context"})
		return
	}
	post, ok := rawPost.(*Post)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid post type in context"})
		return
	}

	// 获取用户ID用于权限校验
	userID, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}
	userIDFloat, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	// 权限校验
	if post.UserID != uint(userIDFloat) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "You are not allowed to delete this post"})
		return
	}

	result := db.Delete(post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
