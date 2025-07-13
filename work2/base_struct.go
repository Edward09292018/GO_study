package work2

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
