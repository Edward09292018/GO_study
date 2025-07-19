package handlers

import (
	"github.com/gin-gonic/gin"
	"go-blog/entity"
	"go-blog/utils"
	"log"
	"net/http"
)

func CreatePost(c *gin.Context) {
	var post entity.Post
	var db = utils.GetDB()
	var postService = NewPostService(db)
	if err := c.BindJSON(&post); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	if err := postService.CreatePost(&post, username.(string)); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}
	log.Print(post.ID)
	SendSuccess(c, post, "Post created successfully", http.StatusOK)

}
func GetPosts(c *gin.Context) {
	var db = utils.GetDB()
	var postService = NewPostService(db)
	posts, err := postService.GetPosts()
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, posts, "Posts retrieved successfully", http.StatusOK)

}
func GetPostByID(c *gin.Context) {
	var db = utils.GetDB()
	var postService = NewPostService(db)
	id := c.Param("id")
	post, err := postService.GetPostByID(id)
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, post, "Post retrieved successfully", http.StatusOK)

}
func UpdatePost(c *gin.Context) {
	var db = utils.GetDB()
	var postService = NewPostService(db)
	id := c.Param("id")
	var updatedPost entity.Post
	if err := c.BindJSON(&updatedPost); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	if err := postService.UpdatePost(id, updatedPost, username.(string)); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, nil, "Post updated successfully", http.StatusOK)

}
func DeletePost(c *gin.Context) {
	var db = utils.GetDB()
	var postService = NewPostService(db)
	id := c.Param("id")
	username, _ := c.Get("username")
	if err := postService.DeletePost(id, username.(string)); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, nil, "Post deleted successfully", http.StatusOK)

}
