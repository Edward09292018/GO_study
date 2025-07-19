package handlers

import (
	"github.com/gin-gonic/gin"
	"go-blog/entity"
	"go-blog/utils"
	"net/http"
)

func CreateComment(c *gin.Context) {
	var db = utils.GetDB()
	var commentService = NewCommentService(db)
	var comment entity.Comment
	if err := c.BindJSON(&comment); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	username, _ := c.Get("username")
	if err := commentService.CreateComment(&comment, username.(string)); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, comment, "Comment created successfully", http.StatusOK)

}
func GetCommentsByPostID(c *gin.Context) {
	var db = utils.GetDB()
	var commentService = NewCommentService(db)
	id := c.Param("id")
	comments, err := commentService.GetCommentsByPostID(id)
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, comments, "Comments retrieved successfully", http.StatusOK)

}
