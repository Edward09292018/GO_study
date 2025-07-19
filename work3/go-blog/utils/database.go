package utils

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectDB() {
	// 从环境变量中读取数据库配置
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	// 使用 GORM 创建数据库连接
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(25)                 // 最大打开的连接数
	sqlDB.SetMaxIdleConns(25)                 // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // 空闲连接最大存活时间

	// 使用 GORM 包装 *sql.DB
	DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 可选：启用 GORM 日志
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
