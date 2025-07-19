package handlers

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(post *Post, username string) error {
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return errors.New("user not found")
	}
	post.UserID = user.ID
	post.CreatedAt = time.Now().Unix()
	post.UpdatedAt = time.Now().Unix()
	if err := s.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (s *PostService) GetPosts() ([]Post, error) {
	var posts []Post
	if err := s.db.Preload("User").Find(&posts).Error; err != nil {
		return nil, errors.New("failed to retrieve posts")
	}
	return posts, nil
}

func (s *PostService) GetPostByID(id string) (Post, error) {
	var post Post
	if err := s.db.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		return post, errors.New("post not found")
	}
	return post, nil
}

func (s *PostService) UpdatePost(id string, updatedPost Post, username string) error {
	var post Post
	if err := s.db.Where("id = ?", id).First(&post).Error; err != nil {
		return errors.New("post not found")
	}

	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if post.UserID != user.ID {
		return errors.New("unauthorized to update this post")
	}

	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	post.UpdatedAt = time.Now().Unix()

	if err := s.db.Save(&post).Error; err != nil {
		return err
	}
	return nil
}

func (s *PostService) DeletePost(id string, username string) error {
	var post Post
	if err := s.db.Where("id = ?", id).First(&post).Error; err != nil {
		return errors.New("post not found")
	}

	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if post.UserID != user.ID {
		return errors.New("unauthorized to delete this post")
	}

	if err := s.db.Delete(&post).Error; err != nil {
		return err
	}
	return nil
}
