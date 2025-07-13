package main

import (
	g6 "github.com/learn/GO_STUDY/work2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func init() {
	dsn := "edward:123@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}
}

// CreatePost 创建文章时更新用户的文章数量
func CreatePost(post *g6.Post, user *g6.User) {
	tx := db.Debug().Model(&g6.User{}).Select("users.id").
		Where("users.id = ?", post.UserID).Find(&user)
	if tx.RowsAffected == 0 {
		panic("用户不存在")
	}
	user.PostCount++ // 自动增加用户文章数
	db.Debug().Save(post)
	db.Debug().Save(user)
}

// DeleteComment 删除评论时检查并更新文章的评论状态
func DeleteComment(comment *g6.Comment, post *g6.Post) {
	tx := db.Debug().Model(&g6.Post{}).Where("id = ?", comment.PostID).Find(&post)
	if tx.RowsAffected == 0 {
		panic("文章不存在")
	}
	post.CommentCount--
	if post.CommentCount == 0 {
		post.HasComments = false
	}
	db.Debug().Delete(comment)
	db.Debug().Save(post)
}
func main() {
	// 示例数据初始化
	user := g6.User{ID: 1, Name: "Alice"}

	post := g6.Post{
		ID: 101, Title: "Go语言入门", UserID: 1, CommentCount: 2, HasComments: true}

	comments := []*g6.Comment{
		{ID: 1, Content: "写得不错！", PostID: 101},
		{ID: 2, Content: "学习了！", PostID: 101},
	}

	//db.AutoMigrate(&g6.User{}, &g6.Post{}, &g6.Comment{})
	//db.Debug().Save(&user)
	//db.Debug().Save(&post)
	//db.Debug().Save(&comments)
	//创建文章并更新用户文章数
	newPost := &g6.Post{ID: 102, Title: "新文章", Content: "内容", UserID: 1}
	CreatePost(newPost, &user)

	// 删除评论并更新文章评论状态
	DeleteComment(comments[0], &post)
	DeleteComment(comments[1], &post)

	// 此时 posts[101].HasComments 应为 false
}
