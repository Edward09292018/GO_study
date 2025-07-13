package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 根据使用的数据库驱动调整
	"github.com/jmoiron/sqlx"
)

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func main() {
	dsn := "edward:123@tcp(localhost:3306)/mydb"
	// 连接数据库
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// 插入一条书籍记录
	_, err = db.NamedExec(
		"INSERT INTO books (title, author, price) VALUES (:title, :author, :price)",
		&Book{
			Title:  "Go Programming",
			Author: "John Doe",
			Price:  69.99,
		},
	)
	if err != nil {
		panic(err)
	}

	// 查询价格大于 50 的书籍
	var books []Book
	err = db.Select(&books, "SELECT id, title, author, price FROM books WHERE price > ?", 50)
	if err != nil {
		panic(err)
	}

	// 输出结果
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Price: %.2f\n", book.ID, book.Title, book.Author, book.Price)
	}

}
