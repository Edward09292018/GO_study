package handlers

import (
	"errors"
	"go-blog/entity"
	"gorm.io/gorm"
	"time"
)

type CommentService struct {
	db *gorm.DB
}
type User = entity.User
type Post = entity.Post
type Comment = entity.Comment

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) CreateComment(comment *Comment, username string) error {
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	var post Post
	if err := s.db.Where("id = ?", comment.PostID).First(&post).Error; err != nil {
		return errors.New("post not found")
	}

	comment.UserID = user.ID
	comment.CreatedAt = time.Now().Unix()

	if err := s.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (s *CommentService) GetCommentsByPostID(postID string) ([]Comment, error) {
	var comments []Comment
	if err := s.db.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, errors.New("failed to retrieve comments")
	}
	return comments, nil
}
