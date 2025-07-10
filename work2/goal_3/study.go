package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

// Employee 员工结构体，对应数据库表字段
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

// GetTechEmployees 查询技术部的所有员工

func GetTechEmployees(db *sqlx.DB) ([]Employee, error) {
	var employees []Employee

	// 使用 Select 执行查询，自动映射到 Employee 结构体切片
	err := db.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")

	if err != nil {
		return nil, fmt.Errorf("failed to get tech employees: %w", err)
	}

	return employees, nil
}

// GetHighestPaidEmployee 查询工资最高的员工

func GetHighestPaidEmployee(db *sqlx.DB) (Employee, error) {
	var employee Employee

	// 使用 Get 执行查询，自动映射到 Employee 结构体
	err := db.Get(&employee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")

	if err != nil {
		return Employee{}, fmt.Errorf("failed to get highest paid employee: %w", err)
	}

	return employee, nil
}
func main() {
	// 数据库连接DSN
	dsn := "edward:123@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 查询技术部所有员工
	techEmployees, err := GetTechEmployees(db)
	if err != nil {
		log.Printf("Error fetching tech employees: %v", err)
	} else {
		fmt.Printf("技术部员工 (%d 位):\n", len(techEmployees))
		for _, emp := range techEmployees {
			fmt.Printf("ID: %d, 姓名: %s, 工资: %.2f\n", emp.ID, emp.Name, emp.Salary)
		}
	}

	// 查询工资最高的员工
	highestPaid, err := GetHighestPaidEmployee(db)
	if err != nil {
		log.Printf("Error fetching highest paid employee: %v", err)
	} else {
		fmt.Printf("\n工资最高的员工:\n")
		fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %.2f\n",
			highestPaid.ID, highestPaid.Name, highestPaid.Department, highestPaid.Salary)
	}
}
