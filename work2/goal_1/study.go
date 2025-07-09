package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Student 定义与数据库表对应的结构体
type Student struct {
	ID    uint   // 自增主键
	Name  string // 学生姓名
	Age   int    // 学生年龄
	Grade string // 学生年级
}

// 2. 创建数据库表
func createTable(db *gorm.DB) {
	err := db.AutoMigrate(&Student{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Table created successfully")
}

// 3. 插入新记录
func insertStudent(db *gorm.DB) {
	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	result := db.Create(&student)

	if result.Error != nil {
		log.Printf("Error inserting record: %v", result.Error)
		return
	}

	log.Printf("Inserted %d records", result.RowsAffected)
}

// 4. 查询年龄大于18岁的学生
func queryStudents(db *gorm.DB) {
	var students []Student
	result := db.Where("age > ?", 18).Find(&students)

	if result.Error != nil {
		log.Printf("Error querying records: %v", result.Error)
		return
	}

	log.Printf("Found %d students older than 18:", result.RowsAffected)
	for _, student := range students {
		log.Printf("ID: %d, Name: %s, Age: %d, Grade: %s",
			student.ID, student.Name, student.Age, student.Grade)
	}
}

// 5. 更新张三年级为四年级
func updateStudentGrade(db *gorm.DB) {
	var student Student
	result := db.Where("name = ?", "张三").First(&student)

	if result.Error != nil {
		log.Printf("Student not found: %v", result.Error)
		return
	}

	student.Grade = "四年级"
	updateResult := db.Save(&student)

	log.Printf("Updated %d records", updateResult.RowsAffected)
}

// 6. 删除年龄小于15岁的学生
func deleteYoungStudents(db *gorm.DB) {
	result := db.Where("age < ?", 15).Delete(&Student{})

	log.Printf("Deleted %d records", result.RowsAffected)
}

// 7. 主函数整合所有操作
func main() {
	// 数据库连接DSN
	dsn := "edward:123@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行各个操作
	createTable(db)
	insertStudent(db)
	queryStudents(db)
	updateStudentGrade(db)
	deleteYoungStudents(db)
}
