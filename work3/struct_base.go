package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Posts    []Post // 用户拥有多篇文章
}
type Post struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Content  string `gorm:"not null"`
	UserID   uint   `gorm:"not null"`
	User     *User
	comments []Comment // 用户拥有多篇文章
}
type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	User    *User
	PostID  uint `gorm:"not null"`
	Post    *Post
}
type ErrorResponse struct {
	Error string `json:"error"`
}

var err error
var db *gorm.DB

func init() {
	dsn := "edward:123@tcp(localhost:3306)/mydb2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

type RequestLog struct {
	gorm.Model
	RequestID string `gorm:"column:request_id" json:"request_id"`
	UserID    string `gorm:"column:user_id" json:"user_id"`
	Method    string `gorm:"column:method" json:"method"`
	Path      string `gorm:"column:path" json:"path"`
	IP        string `gorm:"column:ip" json:"ip"`
	Status    int    `gorm:"column:status" json:"status"`
}

//func main() {
//	err := db.AutoMigrate(&RequestLog{})
//	if err != nil {
//		return
//	}
//}
