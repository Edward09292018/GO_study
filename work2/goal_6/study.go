package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Posts     []Post // 一对多关系：用户有多个文章
	PostCount int
}

type Post struct {
	ID           uint `gorm:"primaryKey"`
	Title        string
	Content      string
	UserID       uint      // 外键，关联 User 表
	User         User      // 关联用户
	Comments     []Comment // 一篇文章可以有多个评论
	CommentCount int
	HasComments  bool
}

type Comment struct {
	ID      uint `gorm:"primaryKey"`
	Content string
	PostID  uint // 外键，关联 Post 表
	Post    Post // 关联文章
}

func main() {
	dsn := "edward:123@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}
	//db.Debug().AutoMigrate(&Post{}, &Comment{}, &User{})

	var user = User{Name: "Edward", Email: "edward@example.com"}

	//使用Gorm查询某个用户发布的所有文章及其对应的评论信息
	//result : user=
	db.Debug().
		Model(&User{}).
		Select("users.*,posts.*,comments.*").
		Joins("left join posts on ?", "posts.user_id = users.id").
		Joins("left join comments on ?", "comments.post_id = posts.id").
		Where(&user).
		Preload("Posts").
		Preload("Posts.Comments").
		Limit(1).
		First(&user)
	fmt.Println(user)
	//- 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	var post Post
	db.Debug().Model(&Post{}).Preload("Comments").
		Select("posts.*").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("count(comments.id) desc").Limit(1).Find(&post)
	fmt.Println(post, len(post.Comments))
}
