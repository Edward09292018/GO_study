package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Posts    []Post
	Comments []Comment
}

type Post struct {
	gorm.Model
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserID    uint
	User      User
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	Comments  []Comment
}

type Comment struct {
	gorm.Model
	Content   string `json:"content"`
	UserID    uint
	User      User
	PostID    uint
	Post      Post
	CreatedAt int64 `json:"created_at"`
}
