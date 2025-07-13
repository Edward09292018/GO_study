package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"unique"`
	Posts []Post // 一对多关系：用户有多个文章
}

type Post struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Content  string
	UserID   uint      // 外键，关联 User 表
	User     User      // 关联用户
	Comments []Comment // 一篇文章可以有多个评论
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

	// 自动迁移模型，创建对应的表
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("数据表迁移失败")
	}

	fmt.Println("数据表创建成功")
}
